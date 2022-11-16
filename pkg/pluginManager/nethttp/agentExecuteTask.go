package nethttp

import (
	"context"
	"fmt"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *PluginNetHttp) AgentEexecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object) (failureReason string, report types.PluginRoundDetail, err error) {
	failureReason = ""
	report = types.PluginRoundDetail{}
	err = nil

	instance, ok := obj.(*crd.Nethttp)
	if !ok {
		msg := "failed to get instance"
		logger.Error(msg)
		err = fmt.Errorf(msg)
		return
	}

	logger.Sugar().Infof("plugin implement task round, instance=%+v", instance)

	return failureReason, report, err
}
