// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package netdns

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	k8sObjManager "github.com/spidernet-io/spiderdoctor/pkg/k8ObjManager"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
	"github.com/spidernet-io/spiderdoctor/pkg/loadRequest"
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"net"
	"strconv"
	"strings"
	"sync"
)

func ParseSuccessCondition(successCondition *crd.NetSuccessCondition, metricResult *loadRequest.DnsMetrics) (failureReason string, err error) {
	switch {
	case successCondition.SuccessRate != nil && metricResult.SuccessRate < *(successCondition.SuccessRate):
		failureReason = fmt.Sprintf("Success Rate %v is lower than request %v", metricResult.SuccessRate, *(successCondition.SuccessRate))
	case successCondition.MeanAccessDelayInMs != nil && metricResult.DelayForSuccess.Mean.Milliseconds() > *(successCondition.MeanAccessDelayInMs):
		failureReason = fmt.Sprintf("mean delay %v ms is bigger than request %v ms", metricResult.DelayForSuccess.Mean.Milliseconds(), *(successCondition.MeanAccessDelayInMs))
	default:
		failureReason = ""
		err = nil
	}
	return
}

func SendRequestAndReport(logger *zap.Logger, targetName string, req *loadRequest.DnsRequestData, successCondition *crd.NetSuccessCondition, report map[string]interface{}) (failureReason string) {

	report["TargetName"] = targetName
	report["TargetServer"] = req.DnsServerAddr
	report["TargetProtocol"] = req.Protocol
	report["Succeed"] = "false"

	result, err := loadRequest.DnsRequest(logger, req)
	if err != nil {
		failureReason = fmt.Sprintf("%v", err)
		logger.Sugar().Errorf("internal error for target %v, error=%v", req.DnsServerAddr, err)
		report["FailureReason"] = failureReason
		return
	}
	report["MeanDelay"] = result.DelayForSuccess.Mean.String()
	report["SucceedRate"] = fmt.Sprintf("%v", result.SuccessRate)

	failureReason, err = ParseSuccessCondition(successCondition, result)
	if err != nil {
		failureReason = fmt.Sprintf("%v", err)
		logger.Sugar().Errorf("internal error for target %v, error=%v", req.DnsServerAddr, err)
		report["FailureReason"] = failureReason
		return
	}

	// generate report
	// notice , upper case for first character of key, or else fail to parse json
	report["Metrics"] = *result
	report["FailureReason"] = failureReason
	if len(report) > 0 {
		report["Succeed"] = "true"
		logger.Sugar().Infof("succeed to test %v", req.DnsServerAddr)
	} else {
		report["Succeed"] = "false"
		logger.Sugar().Warnf("failed to test %v", req.DnsServerAddr)
	}

	return
}

type testTarget struct {
	Name    string
	Request *loadRequest.DnsRequestData
}

func (s *PluginNetDns) AgentExecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object) (finalfailureReason string, finalReport types.PluginRoundDetail, err error) {
	finalfailureReason = ""
	finalReport = types.PluginRoundDetail{}
	err = nil

	instance, ok := obj.(*crd.Netdns)
	if !ok {
		msg := "failed to get instance"
		logger.Error(msg)
		err = fmt.Errorf(msg)
		return
	}

	logger.Sugar().Infof("plugin implement task round, instance=%+v", instance)

	var testTargetList []*testTarget
	var server string

	// Choose whether to request typeA or typeAAAA based on the address type of the server
	if instance.Spec.Target.NetDnsTargetUser != nil {
		server = net.JoinHostPort(*instance.Spec.Target.NetDnsTargetUser.Server, strconv.Itoa(*instance.Spec.Target.NetDnsTargetUser.Port))
		ip := net.ParseIP(*instance.Spec.Target.NetDnsTargetUser.Server)
		if ip.To4() != nil {
			testTargetList = append(testTargetList, &testTarget{Name: "typeA_" + server + "_" + instance.Spec.Request.Domain, Request: &loadRequest.DnsRequestData{
				Protocol:              loadRequest.RequestProtocol(*instance.Spec.Target.Protocol),
				DnsType:               dns.TypeA,
				TargetDomain:          instance.Spec.Request.Domain,
				DnsServerAddr:         server,
				PerRequestTimeoutInMs: int(*instance.Spec.Request.PerRequestTimeoutInMS),
				Qps:                   int(*instance.Spec.Request.QPS),
				DurationInMs:          int(*instance.Spec.Request.DurationInSecond),
			}})
		} else {
			testTargetList = append(testTargetList, &testTarget{Name: "typeAAAA_" + server + "_" + instance.Spec.Request.Domain, Request: &loadRequest.DnsRequestData{
				Protocol:              loadRequest.RequestProtocol(*instance.Spec.Target.Protocol),
				DnsType:               dns.TypeAAAA,
				TargetDomain:          instance.Spec.Request.Domain,
				DnsServerAddr:         server,
				PerRequestTimeoutInMs: int(*instance.Spec.Request.PerRequestTimeoutInMS),
				Qps:                   int(*instance.Spec.Request.QPS),
				DurationInMs:          int(*instance.Spec.Request.DurationInSecond),
			}})
		}
	}

	if instance.Spec.Target.NetDnsTargetDns != nil {
		// When DNS service is not specified, search for DNS services within the cluster
		if instance.Spec.Target.NetDnsTargetDns.ServiceNamespacedName == nil {
			dnsServiceIPs, err := k8sObjManager.GetK8sObjManager().ListServicesDnsIP(ctx)
			if err != nil {
				finalfailureReason = fmt.Sprintf("ListServicesDnsIP err: %v", err)
			}
			logger.Sugar().Infof("dnsServiceIPs %s", dnsServiceIPs)
			for _, serviceIP := range dnsServiceIPs {
				ip := net.ParseIP(serviceIP)
				server = net.JoinHostPort(serviceIP, "53")
				if ip.To4() != nil && *instance.Spec.Target.NetDnsTargetDns.TestIPv4 {
					testTargetList = append(testTargetList, &testTarget{Name: "typeA_" + server + "_" + instance.Spec.Request.Domain, Request: &loadRequest.DnsRequestData{
						Protocol:              loadRequest.RequestProtocol(*instance.Spec.Target.Protocol),
						DnsType:               dns.TypeA,
						TargetDomain:          instance.Spec.Request.Domain,
						DnsServerAddr:         server,
						PerRequestTimeoutInMs: int(*instance.Spec.Request.PerRequestTimeoutInMS),
						Qps:                   int(*instance.Spec.Request.QPS),
						DurationInMs:          int(*instance.Spec.Request.DurationInSecond),
					}})
				} else if ip.To4() == nil && *instance.Spec.Target.NetDnsTargetDns.TestIPv6 {
					testTargetList = append(testTargetList, &testTarget{Name: "typeAAAA_" + server + "_" + instance.Spec.Request.Domain, Request: &loadRequest.DnsRequestData{
						Protocol:              loadRequest.RequestProtocol(*instance.Spec.Target.Protocol),
						DnsType:               dns.TypeAAAA,
						TargetDomain:          instance.Spec.Request.Domain,
						DnsServerAddr:         server,
						PerRequestTimeoutInMs: int(*instance.Spec.Request.PerRequestTimeoutInMS),
						Qps:                   int(*instance.Spec.Request.QPS),
						DurationInMs:          int(*instance.Spec.Request.DurationInSecond),
					}})
				}
			}
		} else {
			// eg: kube-system/coredns
			namespacedName := strings.Split(*instance.Spec.Target.NetDnsTargetDns.ServiceNamespacedName, "/")
			dnsServices, err := k8sObjManager.GetK8sObjManager().GetService(ctx, namespacedName[1], namespacedName[0])
			if err != nil {
				finalfailureReason = fmt.Sprintf("GetService name: %s namespace: %s err: %v", namespacedName[1], namespacedName[0], err)
			}
			for _, serviceIP := range dnsServices.Spec.ClusterIPs {
				ip := net.ParseIP(serviceIP)
				server = net.JoinHostPort(serviceIP, "53")
				if ip.To4() != nil && *instance.Spec.Target.NetDnsTargetDns.TestIPv4 {
					testTargetList = append(testTargetList, &testTarget{Name: "typeA_" + server + "_" + instance.Spec.Request.Domain, Request: &loadRequest.DnsRequestData{
						Protocol:              loadRequest.RequestProtocol(*instance.Spec.Target.Protocol),
						DnsType:               dns.TypeA,
						TargetDomain:          instance.Spec.Request.Domain,
						DnsServerAddr:         server,
						PerRequestTimeoutInMs: int(*instance.Spec.Request.PerRequestTimeoutInMS),
						Qps:                   int(*instance.Spec.Request.QPS),
						DurationInMs:          int(*instance.Spec.Request.DurationInSecond),
					}})
				} else if ip.To4() == nil && *instance.Spec.Target.NetDnsTargetDns.TestIPv6 {
					testTargetList = append(testTargetList, &testTarget{Name: "typeAAAA_" + server + "_" + instance.Spec.Request.Domain, Request: &loadRequest.DnsRequestData{
						Protocol:              loadRequest.RequestProtocol(*instance.Spec.Target.Protocol),
						DnsType:               dns.TypeAAAA,
						TargetDomain:          instance.Spec.Request.Domain,
						DnsServerAddr:         server,
						PerRequestTimeoutInMs: int(*instance.Spec.Request.PerRequestTimeoutInMS),
						Qps:                   int(*instance.Spec.Request.QPS),
						DurationInMs:          int(*instance.Spec.Request.DurationInSecond),
					}})
				}
			}
		}

	}

	var reportList []interface{}

	var wg sync.WaitGroup
	var l lock.Mutex
	for _, item := range testTargetList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, l *lock.Mutex, t testTarget) {
			itemReport := map[string]interface{}{}
			logger.Sugar().Debugf("implement test %v, request %v ", t.Name, *t.Request)
			failureReason := SendRequestAndReport(logger, t.Name, t.Request, instance.Spec.SuccessCondition, itemReport)
			if failureReason != "" {
				finalfailureReason = fmt.Sprintf("test %v: %v", t.Name, failureReason)
			}
			l.Lock()
			reportList = append(reportList, itemReport)
			l.Unlock()
			wg.Done()
		}(&wg, &l, *item)
	}
	wg.Wait()

	logger.Sugar().Infof("plugin finished all http request tests")

	// ----------------------- aggregate report
	finalReport["Detail"] = reportList
	finalReport["TargetType"] = "spiderdoctor agent"
	finalReport["TargetNumber"] = fmt.Sprintf("%d", len(testTargetList))
	if len(finalfailureReason) > 0 {
		logger.Sugar().Errorf("plugin finally failed, %v", finalfailureReason)
		finalReport["FailureReason"] = finalfailureReason
		finalReport["Succeed"] = "false"
	} else {
		finalReport["FailureReason"] = ""
		finalReport["Succeed"] = "true"
	}

	return finalfailureReason, finalReport, err
}
