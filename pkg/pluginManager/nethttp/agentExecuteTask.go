// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nethttp

import (
	"context"
	"fmt"
	k8sObjManager "github.com/spidernet-io/spiderdoctor/pkg/k8ObjManager"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
	"github.com/spidernet-io/spiderdoctor/pkg/loadRequest/loadHttp"
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	config "github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sync"
)

func ParseSuccessCondition(successCondition *crd.NetSuccessCondition, metricResult *loadHttp.Metrics) (failureReason string, err error) {
	switch {
	case successCondition.SuccessRate != nil && float64(metricResult.Success/metricResult.Requests) < *(successCondition.SuccessRate):
		failureReason = fmt.Sprintf("Success Rate %v is lower than request %v", metricResult.Success/metricResult.Requests, *(successCondition.SuccessRate))
	case successCondition.MeanAccessDelayInMs != nil && int64(metricResult.Latencies.Mean) > *(successCondition.MeanAccessDelayInMs):
		failureReason = fmt.Sprintf("mean delay %v ms is bigger than request %v ms", metricResult.Latencies.Mean, *(successCondition.MeanAccessDelayInMs))
	default:
		failureReason = ""
		err = nil
	}
	return
}

func SendRequestAndReport(logger *zap.Logger, targetName string, req *loadHttp.HttpRequestData, successCondition *crd.NetSuccessCondition, report map[string]interface{}) (failureReason string) {

	report["TargetName"] = targetName
	report["TargetUrl"] = req.Url
	report["TargetMethod"] = req.Method
	report["Succeed"] = "false"

	result := loadHttp.HttpRequest(logger, req)
	report["MeanDelay"] = result.Latencies.Mean
	report["SucceedRate"] = fmt.Sprintf("%v", result.Success/result.Requests)

	var err error
	failureReason, err = ParseSuccessCondition(successCondition, result)
	if err != nil {
		failureReason = fmt.Sprintf("%v", err)
		logger.Sugar().Errorf("internal error for target %v, error=%v", req.Url, err)
		report["FailureReason"] = failureReason
		return
	}

	// generate report
	// notice , upper case for first character of key, or else fail to parse json
	report["Metrics"] = *result
	report["FailureReason"] = failureReason
	if len(report) > 0 {
		report["Succeed"] = "true"
		logger.Sugar().Infof("succeed to test %v", req.Url)
	} else {
		report["Succeed"] = "false"
		logger.Sugar().Warnf("failed to test %v", req.Url)
	}

	return
}

type TestTarget struct {
	Name   string
	Url    string
	Method loadHttp.HttpMethod
}

func (s *PluginNetHttp) AgentExecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object) (finalfailureReason string, finalReport types.PluginRoundDetail, err error) {
	finalfailureReason = ""
	finalReport = types.PluginRoundDetail{}
	err = nil
	var e error

	instance, ok := obj.(*crd.Nethttp)
	if !ok {
		msg := "failed to get instance"
		logger.Error(msg)
		err = fmt.Errorf(msg)
		return
	}

	logger.Sugar().Infof("plugin implement task round, instance=%+v", instance)

	target := instance.Spec.Target
	request := instance.Spec.Request
	successCondition := instance.Spec.SuccessCondition
	testTargetList := []*TestTarget{}

	if target.TargetUser != nil {
		logger.Sugar().Infof("load test custom target: Method=%v, Url=%v , qps=%v, PerRequestTimeout=%vs, Duration=%vs", target.TargetUser.Method, target.TargetUser.Url, request.QPS, request.PerRequestTimeoutInMS, request.DurationInSecond)
		finalReport["TargetType"] = "custom url"
		finalReport["TargetNumber"] = "1"
		d := &loadHttp.HttpRequestData{
			Method:              loadHttp.HttpMethod(target.TargetUser.Method),
			Url:                 target.TargetUser.Url,
			Qps:                 request.QPS,
			PerRequestTimeoutMS: request.PerRequestTimeoutInMS,
			RequestTimeSecond:   request.DurationInSecond,
		}
		failureReason := SendRequestAndReport(logger, "custom target", d, successCondition, finalReport)
		if len(failureReason) > 0 {
			finalfailureReason = fmt.Sprintf("test custom target: %v", failureReason)
		}
		return

	} else if target.TargetPod != nil {
		// test pod agent
		logger.Sugar().Infof("load test selected pod: qps=%v, PerRequestTimeout=%vs, Duration=%vs", request.QPS, request.PerRequestTimeoutInMS, request.DurationInSecond)
		finalfailureReason = ""
		var PodIps k8sObjManager.PodIps

		if target.TargetPod.TestMultusInterface {
			PodIps, e = k8sObjManager.GetK8sObjManager().ListSelectedPodMultusIPs(ctx, &(target.TargetPod.PodLabelSelector))
			logger.Sugar().Debugf("test agent multus pod ip: %v", PodIps)
		} else {
			PodIps, e = k8sObjManager.GetK8sObjManager().ListSelectedPodIPs(ctx, &(target.TargetPod.PodLabelSelector))
			logger.Sugar().Debugf("test agent single pod ip: %v", PodIps)
		}
		if e != nil {
			logger.Sugar().Errorf("failed to list por ip, error=%v", e)
			finalfailureReason = fmt.Sprintf("failed to list pod ip, error=%v", e)
		} else if len(PodIps) == 0 {
			s := fmt.Sprintf("failed to find any pod with label %v", target.TargetPod.PodLabelSelector)
			logger.Error(s)
			finalfailureReason = s
		}

		for podname, ips := range PodIps {
			for _, podips := range ips {
				if len(podips.IPv4) > 0 && (target.TargetPod.TestIPv4 == nil || (target.TargetPod.TestIPv4 != nil && *target.TargetPod.TestIPv4)) {
					testTargetList = append(testTargetList, &TestTarget{
						Name:   "SelectedPodV4IP_" + podname + "_" + podips.IPv4,
						Url:    fmt.Sprintf("http://%s:%d%s", podips.IPv4, target.TargetPod.HttpPort, target.TargetPod.UrlPath),
						Method: loadHttp.HttpMethod(target.TargetPod.Method),
					})
				}
				if len(podips.IPv6) > 0 && (target.TargetPod.TestIPv6 == nil || (target.TargetPod.TestIPv6 != nil && *target.TargetPod.TestIPv6)) {
					testTargetList = append(testTargetList, &TestTarget{
						Name:   "SelectedPodV6IP_" + podname + "_" + podips.IPv6,
						Url:    fmt.Sprintf("http://%s:%d%s", podips.IPv6, target.TargetPod.HttpPort, target.TargetPod.UrlPath),
						Method: loadHttp.HttpMethod(target.TargetPod.Method),
					})
				}
			}
		}

	} else {
		// test spiderdoctor agent
		logger.Sugar().Infof("load test spiderdoctor Agent pod: qps=%v, PerRequestTimeout=%vs, Duration=%vs", request.QPS, request.PerRequestTimeoutInMS, request.DurationInSecond)
		finalfailureReason = ""

		// ----------------------- test pod ip
		if target.TargetAgent.TestEndpoint {
			var PodIps k8sObjManager.PodIps

			if target.TargetAgent.TestMultusInterface {
				PodIps, e = k8sObjManager.GetK8sObjManager().ListDaemonsetPodMultusIPs(ctx, config.AgentConfig.Configmap.SpiderDoctorAgentDaemonsetName, config.AgentConfig.PodNamespace)
				logger.Sugar().Debugf("test agent multus pod ip: %v", PodIps)
				if e != nil {
					logger.Sugar().Errorf("failed to ListDaemonsetPodMultusIPs, error=%v", e)
					finalfailureReason = fmt.Sprintf("failed to ListDaemonsetPodMultusIPs, error=%v", e)
				}

			} else {
				PodIps, e = k8sObjManager.GetK8sObjManager().ListDaemonsetPodIPs(ctx, config.AgentConfig.Configmap.SpiderDoctorAgentDaemonsetName, config.AgentConfig.PodNamespace)
				logger.Sugar().Debugf("test agent single pod ip: %v", PodIps)
				if e != nil {
					logger.Sugar().Errorf("failed to ListDaemonsetPodIPs, error=%v", e)
					finalfailureReason = fmt.Sprintf("failed to ListDaemonsetPodIPs, error=%v", e)
				}
			}

			if len(PodIps) > 0 {
				for podname, ips := range PodIps {
					for _, podips := range ips {
						if len(podips.IPv4) > 0 && (target.TargetAgent.TestIPv4 == nil || (target.TargetAgent.TestIPv4 != nil && *target.TargetAgent.TestIPv4)) {
							testTargetList = append(testTargetList, &TestTarget{
								Name:   "AgentPodV4IP_" + podname + "_" + podips.IPv4,
								Url:    fmt.Sprintf("http://%s:%d", podips.IPv4, config.AgentConfig.HttpPort),
								Method: loadHttp.HttpMethodGet,
							})
						}
						if len(podips.IPv6) > 0 && (target.TargetAgent.TestIPv6 == nil || (target.TargetAgent.TestIPv6 != nil && *target.TargetAgent.TestIPv6)) {
							testTargetList = append(testTargetList, &TestTarget{
								Name:   "AgentPodV6IP_" + podname + "_" + podips.IPv6,
								Url:    fmt.Sprintf("http://%s:%d", podips.IPv6, config.AgentConfig.HttpPort),
								Method: loadHttp.HttpMethodGet,
							})
						}
					}
				}
			} else {
				logger.Sugar().Debugf("ignore test agent pod ip")
			}
		}

		// ----------------------- get service
		var agentV4Url, agentV6Url *k8sObjManager.ServiceAccessUrl
		serviceAccessPortName := "http"
		if config.AgentConfig.Configmap.EnableIPv4 {
			agentV4Url, e = k8sObjManager.GetK8sObjManager().GetServiceAccessUrl(ctx, config.AgentConfig.Configmap.AgentSerivceIpv4Name, config.AgentConfig.PodNamespace, serviceAccessPortName)
			if e != nil {
				logger.Sugar().Errorf("failed to get agent ipv4 service url , error=%v", e)
			}
		}
		if config.AgentConfig.Configmap.EnableIPv6 {
			agentV6Url, e = k8sObjManager.GetK8sObjManager().GetServiceAccessUrl(ctx, config.AgentConfig.Configmap.AgentSerivceIpv6Name, config.AgentConfig.PodNamespace, serviceAccessPortName)
			if e != nil {
				logger.Sugar().Errorf("failed to get agent ipv6 service url , error=%v", e)
			}
		}

		var localNodeIpv4, localNodeIpv6 string
		if true {
			localNodeIpv4, localNodeIpv6, e = k8sObjManager.GetK8sObjManager().GetNodeIP(ctx, config.AgentConfig.LocalNodeName)
			if e != nil {
				logger.Sugar().Errorf("failed to get local node %v ip, error=%v", config.AgentConfig.LocalNodeName, e)
			} else {
				logger.Sugar().Debugf("local node %v ip: ipv4=%v, ipv6=%v", config.AgentConfig.LocalNodeName, localNodeIpv4, localNodeIpv6)
			}
		}

		// ----------------------- get ingress
		var agentIngress *networkingv1.Ingress
		agentIngress, e = k8sObjManager.GetK8sObjManager().GetIngress(ctx, config.AgentConfig.Configmap.AgentIngressName, config.AgentConfig.PodNamespace)
		if e != nil {
			logger.Sugar().Errorf("failed to get ingress , error=%v", e)
		}

		// ----------------------- test clusterIP ipv4
		if target.TargetAgent.TestClusterIp && target.TargetAgent.TestIPv4 != nil && *(target.TargetAgent.TestIPv4) {
			if agentV4Url != nil && len(agentV4Url.ClusterIPUrl) > 0 {
				testTargetList = append(testTargetList, &TestTarget{
					Name:   "AgentClusterV4IP_" + agentV4Url.ClusterIPUrl[0],
					Url:    fmt.Sprintf("http://%s", agentV4Url.ClusterIPUrl[0]),
					Method: loadHttp.HttpMethodGet,
				})
			} else {
				finalfailureReason = "failed to get cluster IPv4 IP"
			}
		} else {
			logger.Sugar().Debugf("ignore test agent cluster ipv4 ip")
		}

		// ----------------------- test clusterIP ipv6
		if target.TargetAgent.TestClusterIp && target.TargetAgent.TestIPv6 != nil && *(target.TargetAgent.TestIPv6) {
			if agentV6Url != nil && len(agentV6Url.ClusterIPUrl) > 0 {
				testTargetList = append(testTargetList, &TestTarget{
					Name:   "AgentClusterV6IP_" + agentV6Url.ClusterIPUrl[0],
					Url:    fmt.Sprintf("http://%s", agentV6Url.ClusterIPUrl[0]),
					Method: loadHttp.HttpMethodGet,
				})
			} else {
				finalfailureReason = "failed to get cluster IPv6 IP"
			}
		} else {
			logger.Sugar().Debugf("ignore test agent cluster ipv6 ip")
		}

		// ----------------------- test node port
		if target.TargetAgent.TestNodePort && target.TargetAgent.TestIPv4 != nil && *(target.TargetAgent.TestIPv4) {
			if agentV4Url != nil && agentV4Url.NodePort != 0 && len(localNodeIpv4) != 0 {
				testTargetList = append(testTargetList, &TestTarget{
					Name:   "AgentNodePortV4IP_" + localNodeIpv4 + "_" + fmt.Sprintf("%v", agentV4Url.NodePort),
					Url:    fmt.Sprintf("http://%s:%d", localNodeIpv4, agentV4Url.NodePort),
					Method: loadHttp.HttpMethodGet,
				})
			} else {
				finalfailureReason = "failed to get nodePort IPv4 address"
			}
		} else {
			logger.Sugar().Debugf("ignore test agent nodePort ipv4")
		}

		if target.TargetAgent.TestNodePort && target.TargetAgent.TestIPv6 != nil && *(target.TargetAgent.TestIPv6) {
			if agentV6Url != nil && agentV6Url.NodePort != 0 && len(localNodeIpv6) != 0 {
				testTargetList = append(testTargetList, &TestTarget{
					Name:   "AgentNodePortV6IP_" + localNodeIpv6 + "_" + fmt.Sprintf("%v", agentV6Url.NodePort),
					Url:    fmt.Sprintf("http://%s:%d", localNodeIpv6, agentV6Url.NodePort),
					Method: loadHttp.HttpMethodGet,
				})
			} else {
				finalfailureReason = "failed to get nodePort IPv6 address"
			}
		} else {
			logger.Sugar().Debugf("ignore test agent nodePort ipv6")
		}

		// ----------------------- test loadbalancer IP
		if target.TargetAgent.TestLoadBalancer && target.TargetAgent.TestIPv4 != nil && *(target.TargetAgent.TestIPv4) {
			if agentV4Url != nil && len(agentV4Url.LoadBalancerUrl) > 0 {
				testTargetList = append(testTargetList, &TestTarget{
					Name:   "AgentLoadbalancerV4IP_" + agentV4Url.LoadBalancerUrl[0],
					Url:    fmt.Sprintf("http://%s", agentV4Url.LoadBalancerUrl[0]),
					Method: loadHttp.HttpMethodGet,
				})
			} else {
				finalfailureReason = "failed to get loadbalancer IPv4 address"
			}
		} else {
			logger.Sugar().Debugf("ignore test agent loadbalancer ipv4")
		}

		if target.TargetAgent.TestLoadBalancer && target.TargetAgent.TestIPv6 != nil && *(target.TargetAgent.TestIPv6) {
			if agentV6Url != nil && len(agentV6Url.LoadBalancerUrl) > 0 {
				testTargetList = append(testTargetList, &TestTarget{
					Name:   "AgentLoadbalancerV6IP_" + agentV6Url.LoadBalancerUrl[0],
					Url:    fmt.Sprintf("http://%s", agentV6Url.LoadBalancerUrl[0]),
					Method: loadHttp.HttpMethodGet,
				})
			} else {
				finalfailureReason = "failed to get loadbalancer IPv6 address"
			}
		} else {
			logger.Sugar().Debugf("ignore test agent loadbalancer ipv6")
		}

		// ----------------------- test ingress
		if target.TargetAgent.TestIngress {
			if agentIngress != nil && len(agentIngress.Status.LoadBalancer.Ingress) > 0 {
				http := "http"
				if len(agentIngress.Spec.TLS) > 0 {
					http = "https"
				}
				url := fmt.Sprintf("%s://%s/%s", http, agentIngress.Status.LoadBalancer.Ingress[0].IP, agentIngress.Spec.Rules[0].HTTP.Paths[0].Path)
				testTargetList = append(testTargetList, &TestTarget{
					Name:   "AgentIngress_" + url,
					Url:    url,
					Method: loadHttp.HttpMethodGet,
				})
			} else {
				finalfailureReason = "failed to get agent ingress address"
			}
		} else {
			logger.Sugar().Debugf("ignore test agent ingress ipv6")
		}

	}

	// ------------------------ implement for agent case and selected-pod case
	reportList := []interface{}{}

	var wg sync.WaitGroup
	var l lock.Mutex
	for _, item := range testTargetList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, l *lock.Mutex, t TestTarget) {
			itemReport := map[string]interface{}{}
			d := &loadHttp.HttpRequestData{
				Method:              t.Method,
				Url:                 t.Url,
				Qps:                 request.QPS,
				PerRequestTimeoutMS: request.PerRequestTimeoutInMS,
				RequestTimeSecond:   request.DurationInSecond,
			}
			logger.Sugar().Debugf("implement test %v, request %v ", t.Name, *d)
			failureReason := SendRequestAndReport(logger, t.Name, d, successCondition, itemReport)
			if len(failureReason) > 0 {
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
