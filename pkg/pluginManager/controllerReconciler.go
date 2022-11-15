// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
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

// contorller reconcile
// (1) chedule all task time
// (2) update stauts result
// (3) collect report from agent
func (s *pluginControllerReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	var schedulePlan *crd.SchedulePlan
	var taskStatus *crd.TaskStatus

	if s.plugin.GetApiType().GetDeletionTimestamp() != nil {
		s.logger.Sugar().Debugf("ignore deleting task %v", req)
		return ctrl.Result{}, nil
	}

	handleTask := func(taskStatus *crd.TaskStatus, schedulePlan *crd.SchedulePlan) {
		if taskStatus.ExpectedRound == nil {
			*(taskStatus.ExpectedRound) = schedulePlan.RoundNumber
		}
		if taskStatus.DoneRound == nil {
			*(taskStatus.DoneRound) = 0
		}

		if *taskStatus.DoneRound == *taskStatus.ExpectedRound {
			taskStatus.Finish = true
			s.logger.Sugar().Debugf("ignore finished task %v", req)
			return
		} else {
			taskStatus.Finish = false
		}

	}

	// ------ add crd ------
	switch {
	case s.plugin.GetApiType().GetObjectKind().GroupVersionKind().Kind == KindNameNethttp:
		instance := crd.Nethttp{}
		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		s.logger.Sugar().Debugf("reconcile handle nethttp %v", instance.Name)
		schedulePlan = instance.Spec.Schedule.DeepCopy()
		taskStatus = instance.Status.DeepCopy()

		handleTask(taskStatus, schedulePlan)

	case s.plugin.GetApiType().GetObjectKind().GroupVersionKind().Kind == KindNameNetdns:
		instance := crd.Netdns{}
		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		s.logger.Sugar().Debugf("reconcile handle netdns %v", instance.Name)
		schedulePlan = instance.Spec.Schedule.DeepCopy()
		taskStatus = instance.Status.DeepCopy()

		handleTask(taskStatus, schedulePlan)

	default:
		s.logger.Sugar().Errorf("unknown crd type , detail=%+v", req)
		// forget this
		return ctrl.Result{}, nil
	}

	return s.plugin.ControllerReconcile(s.logger, s.client, ctx, req)
}

var _ reconcile.Reconciler = &pluginControllerReconciler{}

func (s *pluginControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).For(s.plugin.GetApiType()).Complete(s)
}
