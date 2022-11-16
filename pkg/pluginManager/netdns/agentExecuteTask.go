package netdns

import (
	"context"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *PluginNetDns) AgentEexecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object) (result bool, err error) {
	return true, nil

}
