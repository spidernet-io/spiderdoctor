// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package loadRequest

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	"github.com/spidernet-io/spiderdoctor/pkg/utils/stats"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

type RequestProtocol string

const (
	RequestMethodUdp    = RequestProtocol("udp")
	RequestMethodTcp    = RequestProtocol("tcp")
	RequestMethodTcpTls = RequestProtocol("tcp-tls")

	DefaultDnsConfPath = "/etc/resolv.conf"
)

type DnsRequestData struct {
	Protocol RequestProtocol
	// dns.TypeA or dns.TypeAAAA
	DnsType uint16
	// must be full domain
	TargetDomain string
	// empty, or specified to be format "2.2.2.2:53"
	DnsServerAddr         *string
	PerRequestTimeoutInMs int
	Qps                   int
	DurationInMs          int
}

// ------------------

type DelayMetric struct {
	// Mean is the mean request latency.
	Mean time.Duration `json:"mean"`
	// P50 is the 50th percentile request latency.
	P50 time.Duration `json:"50th"`
	// P90 is the 90th percentile request latency.
	P90 time.Duration `json:"90th"`
	// P95 is the 95th percentile request latency.
	P95 time.Duration `json:"95th"`
	// P99 is the 99th percentile request latency.
	P99 time.Duration `json:"99th"`
	// Max is the maximum observed request latency.
	Max time.Duration `json:"max"`
	// Min is the minimum observed request latency.
	Min time.Duration `json:"min"`
}

// final metric
type DnsMetrics struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration

	TargetDomain string
	DnsServer    string
	DnsMethod    string

	// succeed to query the ip
	SucceedCount int
	// failed to get response , or not get ip in the dns response
	FailedCount int
	TotalCount  int
	SuccessRate float64

	// when succeed to get response
	ReplyCode map[string]int
	// error to send request, such as timeout
	ErrorMap map[string]int

	DnsAnswer []dns.RR

	// delay information for success request
	DelayForSuccess DelayMetric
}

// metric for one request
type dnsMetric struct {
	e   error
	rtt time.Duration
	msg *dns.Msg
}

func executeRequestOnce(c *dns.Client, conn *dns.Conn, msg *dns.Msg) *dnsMetric {
	r := dnsMetric{}
	r.msg, r.rtt, r.e = c.ExchangeWithConn(msg, conn)
	return &r
}

func ParseMetrics(final *DnsMetrics, validVals []float32) (*DnsMetrics, error) {
	var e error
	var t float32

	final.SuccessRate = float64(final.SucceedCount) / float64(final.TotalCount)

	// delay
	if final.SucceedCount > 0 {
		t, e = stats.Mean(validVals)
		if e != nil {
			return nil, fmt.Errorf("failed to parse mean delay, error=%v", e)
		}
		final.DelayForSuccess.Mean = time.Duration(t)

		t, e = stats.Max(validVals)
		if e != nil {
			return nil, fmt.Errorf("failed to parse max delay, error=%v", e)
		}
		final.DelayForSuccess.Max = time.Duration(t)

		t, e = stats.Min(validVals)
		if e != nil {
			return nil, fmt.Errorf("failed to parse min delay, error=%v", e)
		}
		final.DelayForSuccess.Min = time.Duration(t)

		t, e = stats.Percentile(validVals, 50)
		if e != nil {
			return nil, fmt.Errorf("failed to parse 50 Percentile, error=%v", e)
		}
		final.DelayForSuccess.P50 = time.Duration(t)

		t, e = stats.Percentile(validVals, 90)
		if e != nil {
			return nil, fmt.Errorf("failed to parse 90 Percentile, error=%v", e)
		}
		final.DelayForSuccess.P90 = time.Duration(t)

		t, e = stats.Percentile(validVals, 95)
		if e != nil {
			return nil, fmt.Errorf("failed to parse 95 Percentile, error=%v", e)
		}
		final.DelayForSuccess.P95 = time.Duration(t)

		t, e = stats.Percentile(validVals, 99)
		if e != nil {
			return nil, fmt.Errorf("failed to parse 99 Percentile, error=%v", e)
		}
		final.DelayForSuccess.P99 = time.Duration(t)
	}

	return final, nil
}

func DnsRequest(logger *zap.Logger, req *DnsRequestData) (result *DnsMetrics, err error) {
	var ServerAddress string
	l := &lock.Mutex{}

	if req.DnsServerAddr == nil {
		config, e := dns.ClientConfigFromFile(DefaultDnsConfPath)
		if e != nil {
			return nil, fmt.Errorf("Error getting nameservers from %v : %v", DefaultDnsConfPath, e)
		}
		if len(config.Servers) == 0 {
			return nil, fmt.Errorf("no name servers in %v ", DefaultDnsConfPath)
		}
		ServerAddress = net.JoinHostPort(config.Servers[0], config.Port)
	} else {
		ServerAddress = *(req.DnsServerAddr)
	}
	// TODO: when dns.TypeAAAA, prefer ipv6 server ?

	logger.Sugar().Infof("dns ServerAddress=%v, request=%v, ", ServerAddress, req)

	if _, ok := dns.IsDomainName(req.TargetDomain); !ok {
		return nil, fmt.Errorf("invalid domain name: %v", req.TargetDomain)
	}
	// if not fqdn, the dns library will report error, so convert the format
	if !dns.IsFqdn(req.TargetDomain) {
		req.TargetDomain = dns.Fqdn(req.TargetDomain)
		logger.Sugar().Debugf("convert target domain to fqdn %v", req.TargetDomain)
	}

	rl := ratelimit.New(req.Qps)
	var wg sync.WaitGroup
	d := time.Duration(req.DurationInMs) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	var duration time.Duration
	logger.Sugar().Infof("begin to request %v for duration %v ", req.TargetDomain, d.String())

	// -------- send all request
	start := time.Now()
	counter := 0

	c := new(dns.Client)
	c.Net = string(req.Protocol)
	c.Timeout = time.Duration(req.PerRequestTimeoutInMs) * time.Millisecond
	msg := new(dns.Msg).SetQuestion(req.TargetDomain, req.DnsType)
	conn, _ := c.Dial(ServerAddress)
	c.SingleInflight = true

	final := &DnsMetrics{
		ErrorMap:  map[string]int{},
		DnsAnswer: []dns.RR{},
		ReplyCode: map[string]int{},
	}

	validVals := []float32{}

	p := func(wg *sync.WaitGroup) {
		r := executeRequestOnce(c, conn, msg)
		l.Lock()
		final.TotalCount++
		if r.e != nil {
			final.FailedCount++
			final.ErrorMap[r.e.Error()]++
		} else {
			if len(r.msg.Answer) > 0 && r.msg.Rcode == dns.RcodeSuccess {
				final.SucceedCount++
				validVals = append(validVals, float32(r.rtt))
			} else {
				final.FailedCount++
			}
			rcodeStr := dns.RcodeToString[r.msg.Rcode]
			final.ReplyCode[rcodeStr]++
		}
		l.Unlock()
		wg.Done()
	}

LOOP:
	for {
		select {
		case <-ctx.Done():
			cancel()
			duration = time.Since(start)
			break LOOP

		default:
			rl.Take()
			counter++
			wg.Add(1)
			go p(&wg)
		}
	}
	wg.Wait()
	end := time.Now()
	logger.Sugar().Infof("finish all %v requests for %v ", counter, req.TargetDomain)
	//-------- parse final metric
	r, e := ParseMetrics(final, validVals)
	if e != nil {
		return nil, fmt.Errorf("failed to parse metric, %v", e)
	}
	r.StartTime = start
	r.EndTime = end
	r.Duration = duration
	r.TargetDomain = req.TargetDomain
	r.DnsServer = ServerAddress
	r.DnsMethod = string(req.Protocol)

	logger.Sugar().Infof("result : %v ", r)
	return r, nil

}
