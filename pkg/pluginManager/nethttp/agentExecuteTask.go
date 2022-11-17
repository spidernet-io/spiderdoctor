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
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
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

func (s *PluginNetHttp) AgentEexecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object) (finalfailureReason string, finalReport types.PluginRoundDetail, err error) {
	finalfailureReason = ""
	finalReport = types.PluginRoundDetail{}
	err = nil

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

	var e error

	// TODO: implement the task
	if target.TargetUrl != nil && len(*target.TargetUrl) != 0 {
		logger.Sugar().Infof("load test custom target: TargetUrl=%v , qps=%v, PerRequestTimeout=%vs, Duration=%vs", *target.TargetUrl, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond)
		finalReport["Type"] = "custom url"
		SendRequestAndReport(logger, *target.TargetUrl, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond, successCondition, finalReport)
		return

	} else {
		// test spiderdoctor agent
		logger.Sugar().Infof("load test spiderdoctor Agent pod: qps=%v, PerRequestTimeout=%vs, Duration=%vs", request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond)

		var finalReport types.PluginRoundDetail
		finalReport["Type"] = "spiderdoctor agent"
		finalfailureReason = ""

		// -- test pod ip
		if target.TargetAgent.TestEndpoint {
			var PodIps k8sObjManager.PodIps
			var report map[string]interface{}

			if target.TargetAgent.TestMultusInterface {
				PodIps, e = k8sObjManager.GetK8sObjManager().ListDaemonsetPodMultusIPs(ctx, config.AgentConfig.SpiderDoctorAgentDaemonsetName, config.AgentConfig.PodNamespace)
				logger.Sugar().Debugf("test agent multus pod ip: %v", PodIps)
				if e != nil {
					logger.Sugar().Errorf("failed to ListDaemonsetPodMultusIPs, error=%v", e)
					report["FailureReason"] = fmt.Sprintf("%v", e)
					finalfailureReason = fmt.Sprintf("%v", e)
				}

			} else {
				PodIps, e = k8sObjManager.GetK8sObjManager().ListDaemonsetPodIPs(ctx, config.AgentConfig.SpiderDoctorAgentDaemonsetName, config.AgentConfig.PodNamespace)
				logger.Sugar().Debugf("test agent single pod ip: %v", PodIps)
				if e != nil {
					logger.Sugar().Errorf("failed to ListDaemonsetPodIPs, error=%v", e)
					report["FailureReason"] = fmt.Sprintf("%v", e)
					finalfailureReason = fmt.Sprintf("%v", e)
				}
			}

			if len(PodIps) > 0 {
				for podname, ips := range PodIps {
					rlist := []interface{}{}
					for _, podips := range ips {
						if len(podips.IPv4) > 0 && (target.TargetAgent.TestIPv4 == nil || (target.TargetAgent.TestIPv4 != nil && *target.TargetAgent.TestIPv4)) {
							var itemReport map[string]interface{}
							target := fmt.Sprintf("http://%s:%d", podips.IPv4, config.AgentConfig.HttpPort)
							logger.Sugar().Debugf("test agent single pod ipv4: %v", target)
							failureReason := SendRequestAndReport(logger, target, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond, successCondition, itemReport)
							if len(failureReason) > 0 {
								finalfailureReason = failureReason
							}
							rlist = append(rlist, itemReport)
						}

						if len(podips.IPv6) > 0 && (target.TargetAgent.TestIPv6 == nil || (target.TargetAgent.TestIPv6 != nil && *target.TargetAgent.TestIPv6)) {
							var itemReport map[string]interface{}
							target := fmt.Sprintf("http://%s:%d", podips.IPv6, config.AgentConfig.HttpPort)
							logger.Sugar().Debugf("test agent single pod ipv6: %v", target)
							failureReason := SendRequestAndReport(logger, target, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond, successCondition, itemReport)
							if len(failureReason) > 0 {
								finalfailureReason = failureReason
							}
							rlist = append(rlist, itemReport)
						}
					}
					report[strings.ToUpper(podname)] = rlist
				}
				finalReport["TestAgentPodIP"] = report
			} else {
				logger.Sugar().Debugf("ignore test agent pod ip: %v")
				finalReport["TestAgentPodIP"] = ""
			}

			// -- test cluster ip

			// -- test node port

			// -- test ingress

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
