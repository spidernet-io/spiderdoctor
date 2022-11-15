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

type pluginAgentReconciler struct {
	client  client.Client
	plugin  plugintypes.ChainingPlugin
	logger  *zap.Logger
	crdKind string
}

var _ reconcile.Reconciler = &pluginAgentReconciler{}

func (s *pluginAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).For(s.plugin.GetApiType()).Complete(s)
}

// https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/internal/controller/controller.go#L320
// when err!=nil , c.Queue.AddRateLimited(req) and log error
// when err==nil && result.Requeue, just c.Queue.AddRateLimited(req)
// when err==nil && result.RequeueAfter > 0 , c.Queue.Forget(obj) and c.Queue.AddAfter(req, result.RequeueAfter)
// or else, c.Queue.Forget(obj)
func (s *pluginAgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	var schedulePlan *crd.SchedulePlan
	var taskStatus *crd.TaskStatus

	if s.plugin.GetApiType().GetDeletionTimestamp() != nil {
		s.logger.Sugar().Debugf("ignore deleting task %v", req)
		return ctrl.Result{}, nil
	}

	handleTask := func(taskStatus *crd.TaskStatus, schedulePlan *crd.SchedulePlan) {

	}

	// ------ add crd ------
	switch s.crdKind {
	case KindNameNethttp:
		instance := crd.Nethttp{}
		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		s.logger.Sugar().Debugf("reconcile handle nethttp %v", instance.Name)
		schedulePlan = instance.Spec.Schedule.DeepCopy()
		taskStatus = instance.Status.DeepCopy()

		handleTask(taskStatus, schedulePlan)

	case KindNameNetdns:
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
		s.logger.Sugar().Errorf("unknown crd type , support kind=%v, detail=%+v", s.crdKind, req)
		// forget this
		return ctrl.Result{}, nil
	}

	return s.plugin.AgentReconcile(s.logger, s.client, ctx, req)
}
