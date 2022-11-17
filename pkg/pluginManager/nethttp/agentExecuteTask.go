package nethttp

import (
	"context"
	"fmt"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/loadRequest"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
)

func ParseSucccessCondition(successCondition *crd.NetSuccessCondition, metricResult *vegeta.Metrics) (failureReason string, err error) {
	switch {
	case metricResult.Success < successCondition.SuccessRate:
		failureReason = fmt.Sprintf("Success Rate %v is lower thant request %v", metricResult.Success, successCondition.SuccessRate)
	case metricResult.Latencies.Mean.Microseconds() > successCondition.MeanAccessDelayInMs:
		failureReason = fmt.Sprintf("mean delay %vs is lower thant request %vs", metricResult.Latencies.Mean.Microseconds(), successCondition.MeanAccessDelayInMs)
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

	// TODO: implement the task
	if target.TargetUrl != nil && len(*target.TargetUrl) != 0 {
		logger.Sugar().Infof("load test custom target: TargetUrl=%v , qps=%v, PerRequestTimeout=%vs, Duration=%vs", *target.TargetUrl, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond)
		result := loadRequest.HttpRequest(*target.TargetUrl, request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond)

		finalfailureReason, err = ParseSucccessCondition(successCondition, result)
		if err != nil {
			logger.Sugar().Errorf("internal error, error=%v", err)
		}
		logger.Sugar().Info("finish, failureReason=%v", finalfailureReason)

		finalReport["target"] = *(target.TargetUrl)
		finalReport["detail"] = *result

		return

	} else {
		// test spiderdoctor agent
		logger.Sugar().Infof("load test spiderdoctor Agent pod: qps=%v, PerRequestTimeout=%vs, Duration=%vs", request.QPS, request.PerRequestTimeoutInSecond, request.DurationInSecond)

	}

	return finalfailureReason, finalReport, err
}
