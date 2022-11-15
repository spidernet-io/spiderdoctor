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
	plugin plugintypes.ChainingPlugin
	logger *zap.Logger
}

var _ reconcile.Reconciler = &pluginAgentReconciler{}

func (s *pluginAgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return s.plugin.AgentReconcile(s.logger, s.client, ctx, req)
}

func (s *pluginAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(s.plugin.GetApiType()).Complete(s)
}
