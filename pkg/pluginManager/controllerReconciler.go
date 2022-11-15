package pluginManager

import (
	"context"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type pluginControllerReconciler struct {
	client client.Client
	plugin plugintypes.ChainingPlugin
	logger *zap.Logger
}

func (s *pluginControllerReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	return s.plugin.ControllerReconcile(s.logger, s.client, ctx, r)
}

var _ reconcile.Reconciler = &pluginControllerReconciler{}

func (s *pluginControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(s.plugin.GetApiType()).Complete(s)
}
