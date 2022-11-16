package podManager

import (
	"context"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PodManager interface {
	SetupWithManager(mgr ctrl.Manager) error
}

type podManager struct {
	client client.Client
	logger *zap.Logger
}

func NewManager(client client.Client, logger *zap.Logger) PodManager {
	return &podManager{
		client: client,
		logger: logger,
	}
}
func (s *podManager) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func (s *podManager) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&corev1.Pod{}).Complete(s)
}
