package pluginManager

import (
	"context"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type pluginAgentReconciler struct {
	client client.Client
	p      plugintypes.ChainingPlugin
	logger *zap.Logger
}

func (s *pluginAgentReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	return s.p.AgentReconcile(s.logger, s.client, ctx, r)
}

var _ reconcile.Reconciler = &pluginAgentReconciler{}

func (s *pluginManager) runAgentReconcile() {
	logger := s.logger

	n := ctrl.Options{
		MetricsBindAddress:     "0",
		HealthProbeBindAddress: "0",
		LeaderElection:         false,
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), n)
	if err != nil {
		logger.Sugar().Fatalf("failed to NewManager, reason=%v", err)
	}
	builder := ctrl.NewControllerManagedBy(mgr)
	for name, plugin := range s.chainingPlugins {
		go func(name string, t plugintypes.ChainingPlugin) {
			logger.Sugar().Infof("run controller for plugin %v", name)
			builder.For(t.GetApiType()).Owns(t.GetApiType()).Build(&pluginAgentReconciler{logger: logger.Named(name + "Reconciler"), p: t})
		}(name, plugin)
	}
}
