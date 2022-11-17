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
		var agentV4Service *corev1.Service
		var agentV6Service *corev1.Service
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
			if config.AgentConfig.Configmap.EnableIPv4 {
				agentV4Service, e = k8sObjManager.GetK8sObjManager().GetService(ctx, config.AgentConfig.AgentSerivceIpv4Name, config.AgentConfig.PodNamespace)
				if e != nil {
					logger.Sugar().Errorf("failed to get agent ipv4 service, error=%v", e)
				}
			}
			if config.AgentConfig.Configmap.EnableIPv6 {
				agentV6Service, e = k8sObjManager.GetK8sObjManager().GetService(ctx, config.AgentConfig.AgentSerivceIpv6Name, config.AgentConfig.PodNamespace)
				if e != nil {
					logger.Sugar().Errorf("failed to get agent ipv6 service, error=%v", e)
				}
			}

			// ----------------------- test cluster ipv4 ip
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

			// ----------------------- test cluster ipv6 ip
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

			// ----------------------- test loadbalancer IP

			// ----------------------- test ingress

			reportList := []interface{}{}
			for _, targetItem := range testTargetList {
				itemReport := map[string]interface{}{}
				logger.Sugar().Debugf("implement test %v, target=%v", targetItem.Name, targetItem.Url)
				failureReason := SendRequestAndReport(logger, targetItem.Name, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond, successCondition, itemReport)
				if len(failureReason) > 0 {
					finalfailureReason = fmt.Sprintf("test %v: %v", targetItem.Name, failureReason)
				}
				reportList = append(reportList, itemReport)
			}

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
