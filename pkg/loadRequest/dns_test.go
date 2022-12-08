// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package loadRequest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/spiderdoctor/pkg/loadRequest"
	"github.com/spidernet-io/spiderdoctor/pkg/logger"
)

var _ = Describe("test dns ", Label("dns"), func() {

	It("test udp ", func() {

		dnsServer := "223.5.5.5:53"
		req := &loadRequest.DnsRequestData{
			Protocol:              loadRequest.RequestMethodUdp,
			DnsType:               dns.TypeA,
			TargetDomain:          "www.baidu.com",
			DnsServerAddr:         &dnsServer,
			PerRequestTimeoutInMs: 1000,
			DurationInMs:          1000,
			Qps:                   10,
		}

		log := logger.NewStdoutLogger("debug", "test")
		result, e := loadRequest.DnsRequest(log, req)
		Expect(e).NotTo(HaveOccurred(), "failed to execute , error=%v", e)
		Expect(result.FailedCount).To(Equal(0))
		Expect(len(result.ReplyCode)).To(Equal(1))
		Expect(result.ReplyCode).Should(HaveKey(dns.RcodeToString[dns.RcodeSuccess]))

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())

	})

	It("test tcp ", func() {

		dnsServer := "223.5.5.5:53"
		req := &loadRequest.DnsRequestData{
			Protocol:              loadRequest.RequestMethodTcp,
			DnsType:               dns.TypeA,
			TargetDomain:          "www.baidu.com",
			DnsServerAddr:         &dnsServer,
			PerRequestTimeoutInMs: 1000,
			DurationInMs:          1000,
			Qps:                   10,
		}

		log := logger.NewStdoutLogger("debug", "test")
		result, e := loadRequest.DnsRequest(log, req)
		Expect(e).NotTo(HaveOccurred(), "failed to execute , error=%v", e)
		Expect(result.FailedCount).To(Equal(0))
		Expect(len(result.ReplyCode)).To(Equal(1))
		Expect(result.ReplyCode).Should(HaveKey(dns.RcodeToString[dns.RcodeSuccess]))

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())
	})

	It("test bad domain ", func() {

		dnsServer := "223.5.5.5:53"
		req := &loadRequest.DnsRequestData{
			Protocol:              loadRequest.RequestMethodUdp,
			DnsType:               dns.TypeA,
			TargetDomain:          "www.no-existed.com",
			DnsServerAddr:         &dnsServer,
			PerRequestTimeoutInMs: 1000,
			DurationInMs:          1000,
			Qps:                   10,
		}

		log := logger.NewStdoutLogger("debug", "test")
		result, e := loadRequest.DnsRequest(log, req)
		Expect(e).NotTo(HaveOccurred(), "failed to execute , error=%v", e)
		Expect(result.SucceedCount).To(Equal(0))
		Expect(len(result.ReplyCode)).To(Equal(1))
		Expect(result.ReplyCode).Should(HaveKey(dns.RcodeToString[dns.RcodeNameError]))

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())

	})

	It("test aaaa ", Label("aaaa"), func() {
		dnsServer := "223.5.5.5:53"
		req := &loadRequest.DnsRequestData{
			Protocol:              loadRequest.RequestMethodUdp,
			DnsType:               dns.TypeAAAA,
			TargetDomain:          "wikipedia.org",
			DnsServerAddr:         &dnsServer,
			PerRequestTimeoutInMs: 1000,
			DurationInMs:          1000,
			Qps:                   10,
		}

		log := logger.NewStdoutLogger("debug", "test")
		result, e := loadRequest.DnsRequest(log, req)
		Expect(e).NotTo(HaveOccurred(), "failed to execute , error=%v", e)
		Expect(result.FailedCount).To(Equal(0))
		Expect(len(result.ReplyCode)).To(Equal(1))
		Expect(result.ReplyCode).Should(HaveKey(dns.RcodeToString[dns.RcodeSuccess]))

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())

	})
})
