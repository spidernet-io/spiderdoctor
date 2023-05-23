// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package loadDns

import (
	"github.com/miekg/dns"
	"github.com/spidernet-io/spiderdoctor/pkg/utils/stats"
	"sync"
	"time"
)

// Max size of the buffer of result channel.
const maxResult = 1000000

type result struct {
	err      error
	duration time.Duration
	msg      *dns.Msg
}

type Metrics struct {
	StartTime    time.Time           `json:"start"`
	EndTime      time.Time           `json:"end"`
	TargetDomain string              `json:"target_domain"`
	DnsServer    string              `json:"dns_server"`
	DnsMethod    string              `json:"Method"`
	Duration     string              `json:"duration"`
	Requests     int64               `json:"requests"`
	Success      int64               `json:"success"`
	Failed       int64               `json:"failed"`
	TPS          float64             `json:"TPS"`
	Latencies    latencyDistribution `json:"latencies"`
	Errors       map[string]int      `json:"errors"`
	ReplyCode    map[string]int      `json:"reply_code"`
}

type latencyDistribution struct {
	// Mean is the mean request latency.
	Mean float32 `json:"Mean_inMs"`
	// P50 is the 50th percentile request latency.
	P50 float32 `json:"P50_inMs"`
	// P90 is the 90th percentile request latency.
	P90 float32 `json:"P90_inMs"`
	// P95 is the 95th percentile request latency.
	P95 float32 `json:"P95_inMs"`
	// P99 is the 99th percentile request latency.
	P99 float32 `json:"P99_inMs"`
	// Max is the maximum observed request latency.
	Max float32 `json:"Max_inMs"`
	// Min is the minimum observed request latency.
	Min float32 `json:"Min_inMs"`
}

type Work struct {
	ServerAddr string

	Msg *dns.Msg

	Protocol string

	// C is the concurrency level, the number of concurrent workers to run.
	Concurrency int

	// Timeout in seconds.
	Timeout int

	// Qps is the rate limit in queries per second.
	QPS int

	initOnce  sync.Once
	results   chan *result
	stopCh    chan struct{}
	start     time.Duration
	startTime time.Time
	report    *report
}

// Init initializes internal data-structures
func (b *Work) Init() {
	b.initOnce.Do(func() {
		b.results = make(chan *result, maxResult)
		b.stopCh = make(chan struct{}, b.Concurrency)
	})
}

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *Work) Run() {
	b.Init()
	b.startTime = time.Now()
	b.start = time.Since(b.startTime)
	b.report = newReport(b.results)
	// Run the reporter first, it polls the result channel until it is closed.
	go func() {
		runReporter(b.report)
	}()

	b.runWorkers()
	b.Finish()
}

func (b *Work) Stop() {
	// Send stop signal so that workers can stop gracefully.
	for i := 0; i < b.Concurrency; i++ {
		b.stopCh <- struct{}{}
	}
}

func (b *Work) Finish() {
	close(b.results)
	total := b.now() - b.start
	// Wait until the reporter is done.
	<-b.report.done
	b.report.finalize(total)
}

func (b *Work) makeRequest(client *dns.Client, conn *dns.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	msg, rtt, err := client.ExchangeWithConn(b.Msg, conn)
	b.results <- &result{
		duration: rtt,
		err:      err,
		msg:      msg,
	}
}

func (b *Work) runWorker() {
	var ticker *time.Ticker
	if b.QPS > 0 {
		ticker = time.NewTicker(time.Duration(1e6*b.Concurrency/(b.QPS)) * time.Microsecond)
	}
	client := new(dns.Client)
	client.Net = b.Protocol
	client.Timeout = time.Duration(b.Timeout) * time.Millisecond
	conn, _ := client.Dial(b.ServerAddr)
	client.SingleInflight = true
	wg := &sync.WaitGroup{}
	for {
		// Check if application is stopped. Do not send into a closed channel.
		select {
		case <-b.stopCh:
			wg.Wait()
			return
		default:
			if b.QPS > 0 {
				<-ticker.C
			}
			wg.Add(1)

			// check connect close
			// if close new connect
			if conn == nil {
				conn, _ = client.Dial(b.ServerAddr)
			}
			go b.makeRequest(client, conn, wg)
		}
	}
}

func (b *Work) runWorkers() {
	var wg sync.WaitGroup
	wg.Add(b.Concurrency)
	for i := 0; i < b.Concurrency; i++ {
		go func() {
			b.runWorker()
			wg.Done()
		}()
	}
	wg.Wait()

}

func (b *Work) now() time.Duration { return time.Since(b.startTime) }

func (b *Work) AggregateMetric() *Metrics {
	metric := &Metrics{}
	metric.Requests = b.report.totalCount
	metric.StartTime = b.startTime
	metric.EndTime = b.startTime.Add(b.report.total)
	metric.Duration = b.report.total.String()
	metric.Failed = b.report.failedCount
	metric.Success = b.report.successCount
	metric.TPS = b.report.tps
	latency := latencyDistribution{}

	t, _ := stats.Mean(b.report.lats)
	latency.Mean = t

	t, _ = stats.Max(b.report.lats)
	latency.Max = t

	t, _ = stats.Min(b.report.lats)

	latency.Min = t

	t, _ = stats.Percentile(b.report.lats, 50)
	latency.P50 = t

	t, _ = stats.Percentile(b.report.lats, 90)
	latency.P90 = t

	t, _ = stats.Percentile(b.report.lats, 95)
	latency.P95 = t

	t, _ = stats.Percentile(b.report.lats, 99)
	latency.P99 = t

	metric.Latencies = latency

	metric.Errors = b.report.errorDist

	metric.ReplyCode = b.report.ReplyCode

	metric.TargetDomain = b.Msg.Question[0].Name

	metric.DnsMethod = b.Protocol

	metric.DnsServer = b.ServerAddr

	return metric
}
