package nethttp

import (
	"context"
	"fmt"
	k8sObjManager "github.com/spidernet-io/spiderdoctor/pkg/k8ObjManager"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/loadRequest"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	config "github.com/spidernet-io/spiderdoctor/pkg/types"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

func ParseSucccessCondition(successCondition *crd.NetSuccessCondition, metricResult *vegeta.Metrics) (failureReason string, err error) {
	switch {
	case metricResult.Success < successCondition.SuccessRate:
		failureReason = fmt.Sprintf("Success Rate %v is lower thant request %v", metricResult.Success, successCondition.SuccessRate)
	case metricResult.Latencies.Mean.Microseconds() > successCondition.MeanAccessDelayInMs:
		failureReason = fmt.Sprintf("mean delay %vs is lower thant request %vs", metricResult.Latencies.Mean.Microseconds(), successCondition.MeanAccessDelayInMs)
	default:
		failureReason = ""
		err = nil
	}
	return
}

func SendRequestAndReport(logger *zap.Logger, TargetUrl string, qps, PerRequestTimeoutInSecond, DurationInSecond int, successCondition *crd.NetSuccessCondition, report map[string]interface{}) (failureReason string) {
	failureReason = ""

	report["Target"] = TargetUrl
	report["Succeed"] = "false"

	result := loadRequest.HttpRequest(TargetUrl, qps, PerRequestTimeoutInSecond, DurationInSecond)

	var err error
	failureReason, err = ParseSucccessCondition(successCondition, result)
	if err != nil {
		failureReason = fmt.Sprintf("%v", err)
		logger.Sugar().Errorf("internal error for target %v, error=%v", TargetUrl, err)
		report["FailureReason"] = failureReason
		return
	}

	// generate report
	// notice , upper case for first character of key, or else fail to parse json
	report["Detail"] = *result
	report["FailureReason"] = failureReason
	if len(report) > 0 {
		report["Succeed"] = "true"
		logger.Sugar().Infof("succeed to test %v", TargetUrl)
	} else {
		report["Succeed"] = "false"
		logger.Sugar().Warnf("failed to test %v", TargetUrl)
	}

	return
}

type TestTarget struct {
	Name string
	Url  string
}

func (s *PluginNetHttp) AgentEexecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object) (finalfailureReason string, finalReport types.PluginRoundDetail, err error) {
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

	plan := instance.Spec.Schedule
	target := instance.Spec.Target
	request := instance.Spec.Request
	successCondition := instance.Spec.SuccessCondition

	if target.TargetUrl != nil && len(*target.TargetUrl) != 0 {
		logger.Sugar().Infof("load test custom target: TargetUrl=%v , qps=%v, PerRequestTimeout=%vs, Duration=%vs", *target.TargetUrl, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond)
		finalReport["Type"] = "custom url"
		SendRequestAndReport(logger, *target.TargetUrl, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond, successCondition, finalReport)
		return

	} else {
		// test spiderdoctor agent
		logger.Sugar().Infof("load test spiderdoctor Agent pod: qps=%v, PerRequestTimeout=%vs, Duration=%vs", request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond)
		finalfailureReason = ""
		testTargetList := []*TestTarget{}

		// ----------------------- test pod ip
		if target.TargetAgent.TestEndpoint {
			var PodIps k8sObjManager.PodIps

			if target.TargetAgent.TestMultusInterface {
				PodIps, e = k8sObjManager.GetK8sObjManager().ListDaemonsetPodMultusIPs(ctx, config.AgentConfig.SpiderDoctorAgentDaemonsetName, config.AgentConfig.PodNamespace)
				logger.Sugar().Debugf("test agent multus pod ip: %v", PodIps)
				if e != nil {
					logger.Sugar().Errorf("failed to ListDaemonsetPodMultusIPs, error=%v", e)
					finalfailureReason = fmt.Sprintf("failed to ListDaemonsetPodMultusIPs, error=%v", e)
				}

			} else {
				PodIps, e = k8sObjManager.GetK8sObjManager().ListDaemonsetPodIPs(ctx, config.AgentConfig.SpiderDoctorAgentDaemonsetName, config.AgentConfig.PodNamespace)
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
								Name: "AgentPodV4IP_" + podname + "_" + podips.IPv4,
								Url:  fmt.Sprintf("http://%s:%d", podips.IPv4, config.AgentConfig.HttpPort),
							})
						}
						if len(podips.IPv6) > 0 && (target.TargetAgent.TestIPv6 == nil || (target.TargetAgent.TestIPv6 != nil && *target.TargetAgent.TestIPv6)) {
							testTargetList = append(testTargetList, &TestTarget{
								Name: "AgentPodV6IP_" + podname + "_" + podips.IPv6,
								Url:  fmt.Sprintf("http://%s:%d", podips.IPv6, config.AgentConfig.HttpPort),
							})
						}
					}
				}
			} else {
				logger.Sugar().Debugf("ignore test agent pod ip")
			}

			// ----------------------- get service
			var serviceNodePortV4, serviceNodePortV6 int32
			var agentV4Service *corev1.Service
			var agentV6Service *corev1.Service
			var localNodeIpv4, localNodeIpv6 string

			if config.AgentConfig.Configmap.EnableIPv4 {
				agentV4Service, e = k8sObjManager.GetK8sObjManager().GetService(ctx, config.AgentConfig.AgentSerivceIpv4Name, config.AgentConfig.PodNamespace)
				if e != nil {
					logger.Sugar().Errorf("failed to get agent ipv4 service, error=%v", e)
				} else {
					logger.Sugar().Debugf("agent ipv4 service: %v", agentV4Service.Spec)
					// find nodePort
					for _, v := range agentV4Service.Spec.Ports {
						if v.Name == "http" && v.NodePort != 0 {
							serviceNodePortV4 = v.NodePort
							logger.Sugar().Debugf("agent ipv4 service nodePort: %v", serviceNodePortV4)
							break
						}
					}
				}
			}
			if config.AgentConfig.Configmap.EnableIPv6 {
				agentV6Service, e = k8sObjManager.GetK8sObjManager().GetService(ctx, config.AgentConfig.AgentSerivceIpv6Name, config.AgentConfig.PodNamespace)
				if e != nil {
					logger.Sugar().Errorf("failed to get agent ipv6 service, error=%v", e)
				} else {
					logger.Sugar().Debugf("agent ipv6 service: %v", agentV6Service.Spec)
					// find nodePort
					for _, v := range agentV6Service.Spec.Ports {
						if v.Name == "http" && v.NodePort != 0 {
							serviceNodePortV6 = v.NodePort
							logger.Sugar().Debugf("agent ipv6 service nodePort: %v", serviceNodePortV6)
							break
						}
					}
				}
			}
			if true {
				localNodeIpv4, localNodeIpv6, e = k8sObjManager.GetK8sObjManager().GetNodeIP(ctx, config.AgentConfig.LocalNodeName)
				if e != nil {
					logger.Sugar().Errorf("failed to get local node %v ip, error=%v", config.AgentConfig.LocalNodeName, e)
				} else {
					logger.Sugar().Debugf("local node %v ip: ipv4=%v, ipv6=%v", config.AgentConfig.LocalNodeName, localNodeIpv4, localNodeIpv6)
				}
			}

			// ----------------------- test clusterIP ipv4
			if target.TargetAgent.TestClusterIp && target.TargetAgent.TestIPv4 != nil && *(target.TargetAgent.TestIPv4) {
				if agentV4Service != nil && len(agentV4Service.Spec.ClusterIP) != 0 {
					testTargetList = append(testTargetList, &TestTarget{
						Name: "AgentClusterV4IP_" + agentV4Service.Spec.ClusterIP,
						Url:  fmt.Sprintf("http://%s:%d", agentV4Service.Spec.ClusterIP, config.AgentConfig.HttpPort),
					})
				} else {
					finalfailureReason = "failed to get cluster IPv4 IP"
				}
			} else {
				logger.Sugar().Debugf("ignore test agent cluster ipv4 ip")
			}

			// ----------------------- test clusterIP ipv6
			if target.TargetAgent.TestClusterIp && target.TargetAgent.TestIPv6 != nil && *(target.TargetAgent.TestIPv6) {
				reportRoot := map[string]interface{}{}
				if agentV6Service == nil {
					testTargetList = append(testTargetList, &TestTarget{
						Name: "AgentClusterV6IP_" + agentV6Service.Spec.ClusterIP,
						Url:  fmt.Sprintf("http://%s:%d", agentV6Service.Spec.ClusterIP, config.AgentConfig.HttpPort),
					})
				} else {
					finalfailureReason = "failed to get cluster IPv6 IP"
				}
				finalReport["TestAgentClusterIPv6IP"] = reportRoot
			} else {
				logger.Sugar().Debugf("ignore test agent cluster ipv6 ip")
			}

			// ----------------------- test node port
			if target.TargetAgent.TestNodePort && target.TargetAgent.TestIPv4 != nil && *(target.TargetAgent.TestIPv4) {
				if agentV4Service != nil && len(localNodeIpv4) != 0 && serviceNodePortV4 != 0 {
					testTargetList = append(testTargetList, &TestTarget{
						Name: "AgentNodePortV4IP_" + localNodeIpv4 + "_" + fmt.Sprintf("%v", serviceNodePortV4),
						Url:  fmt.Sprintf("http://%s:%d", localNodeIpv4, serviceNodePortV4),
					})
				} else {
					finalfailureReason = "failed to get nodePort IPv4 address"
				}
			} else {
				logger.Sugar().Debugf("ignore test agent nodePort ipv4")
			}
			if target.TargetAgent.TestNodePort && target.TargetAgent.TestIPv6 != nil && *(target.TargetAgent.TestIPv6) {
				if agentV6Service != nil && len(localNodeIpv6) != 0 && serviceNodePortV6 != 0 {
					testTargetList = append(testTargetList, &TestTarget{
						Name: "AgentNodePortV6IP_" + localNodeIpv6 + "_" + fmt.Sprintf("%v", serviceNodePortV6),
						Url:  fmt.Sprintf("http://%s:%d", localNodeIpv6, serviceNodePortV6),
					})
				} else {
					finalfailureReason = "failed to get nodePort IPv6 address"
				}
			} else {
				logger.Sugar().Debugf("ignore test agent nodePort ipv6")
			}

			// TODO: ----------------------- test loadbalancer IP

			// TODO: ----------------------- test ingress

			// ------------------------ implement it
			reportList := []interface{}{}
			testNum := len(testTargetList)
			if testNum*request.DurationInSecond < (int(plan.TimeoutMinute) * 60) {
				logger.Sugar().Infof("plugin implement %v tests, it takes about %vs, shorter than required %vs ", testNum, testNum*request.DurationInSecond, plan.TimeoutMinute*60)
			} else {
				logger.Sugar().Errorf("plugin implement %v tests, it takes about %vs, logger than required %vs ", testNum, testNum*request.DurationInSecond, plan.TimeoutMinute*60)
			}
			start := time.Now()
			for _, targetItem := range testTargetList {
				itemReport := map[string]interface{}{}
				logger.Sugar().Debugf("implement test %v, target=%v , QPS=%v, PerRequestTimeoutInSecond=%v, DurationInSecond=%v ", targetItem.Name, targetItem.Url, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond)
				failureReason := SendRequestAndReport(logger, targetItem.Url, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond, successCondition, itemReport)
				if len(failureReason) > 0 {
					finalfailureReason = fmt.Sprintf("test %v: %v", targetItem.Name, failureReason)
				}
				reportList = append(reportList, itemReport)
			}
			logger.Sugar().Infof("plugin finished %v tests, it taked time with %v , started at %v", testNum, testNum*request.DurationInSecond, time.Now().Sub(start).String())

			// ----------------------- aggregate report
			finalReport["Detail"] = reportList
			finalReport["Type"] = "spiderdoctor agent"
			if len(finalfailureReason) > 0 {
				finalReport["FailureReason"] = finalfailureReason
				finalReport["Succeed"] = "false"
			} else {
				finalReport["FailureReason"] = ""
				finalReport["Succeed"] = "true"
			}
		}

		return finalfailureReason, finalReport, err
	}
}
