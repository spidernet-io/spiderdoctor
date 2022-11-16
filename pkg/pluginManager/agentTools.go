package pluginManager

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (s *pluginAgentReconciler) HandleAgentTaskRound(logger *zap.Logger, ctx context.Context, oldStatus *crd.TaskStatus, schedulePlan *crd.SchedulePlan, taskName string) (result *reconcile.Result, taskStatus *crd.TaskStatus, e error) {
	newStatus := oldStatus.DeepCopy()
	return result, newStatus, nil
}
