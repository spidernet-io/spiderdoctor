// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package httpapphealthy

import (
	"context"
	"fmt"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
	"github.com/spidernet-io/spiderdoctor/pkg/loadRequest/loadHttp"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
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

func (s *PluginHttpAppHealthy) AgentExecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object) (finalfailureReason string, finalReport types.PluginRoundDetail, err error) {
	finalfailureReason = ""
	finalReport = types.PluginRoundDetail{}
	err = nil

	instance, ok := obj.(*crd.HttpAppHealthy)
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

	logger.Sugar().Infof("load test custom target: Method=%v, Url=%v , qps=%v, PerRequestTimeout=%vs, Duration=%vs", target.Method, target.Host, request.QPS, request.PerRequestTimeoutInMS, request.DurationInSecond)
	finalReport["TargetType"] = "HttpAppHealthy"
	finalReport["TargetNumber"] = "1"
	d := &loadHttp.HttpRequestData{
		Method:              loadHttp.HttpMethod(target.Method),
		Url:                 target.Host,
		Qps:                 request.QPS,
		PerRequestTimeoutMS: request.PerRequestTimeoutInMS,
		RequestTimeSecond:   request.DurationInSecond,
		Http2:               target.Http2,
	}
	failureReason := SendRequestAndReport(logger, "HttpAppHealthy target", d, successCondition, finalReport)
	if len(failureReason) > 0 {
		finalfailureReason = fmt.Sprintf("test HttpAppHealthy target: %v", failureReason)
	}

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
