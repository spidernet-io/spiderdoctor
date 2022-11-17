// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"github.com/spidernet-io/spiderdoctor/pkg/taskStatusManager"
	"go.uber.org/zap"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type pluginAgentReconciler struct {
	client        client.Client
	plugin        plugintypes.ChainingPlugin
	logger        *zap.Logger
	crdKind       string
	localNodeName string
	taskRoundData taskStatusManager.TaskStatus
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

	// ------ add crd ------
	switch s.crdKind {
	case KindNameNethttp:
		instance := crd.Nethttp{}
		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		logger := s.logger.With(zap.String(instance.Kind, instance.Name))
		logger.Sugar().Debugf("reconcile handle %v", instance)
		if instance.DeletionTimestamp != nil {
			s.logger.Sugar().Debugf("ignore deleting task %v", req)
			return ctrl.Result{}, nil
		}

		oldStatus := instance.Status.DeepCopy()
		taskName := instance.Kind + "." + instance.Name
		if result, newStatus, err := s.HandleAgentTaskRound(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), &instance, taskName); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to HandleAgentTaskRound, will retry it, error=%v", err)
			return ctrl.Result{}, err

		} else {
			if newStatus != nil && !reflect.DeepEqual(newStatus, oldStatus) {
				instance.Status = *newStatus
				if err := s.client.Status().Update(ctx, &instance); err != nil {
					// requeue
					logger.Sugar().Errorf("failed to update status, will retry it, error=%v", err)
					return ctrl.Result{}, err
				}
				logger.Sugar().Debugf("succeeded update status, newStatus=%+v", newStatus)
			}

			if result != nil {
				return *result, nil
			}
		}

	case KindNameNetdns:
		instance := crd.Netdns{}
		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		logger := s.logger.With(zap.String(instance.Kind, instance.Name))
		logger.Sugar().Debugf("reconcile handle %v", instance)
		if instance.DeletionTimestamp != nil {
			s.logger.Sugar().Debugf("ignore deleting task %v", req)
			return ctrl.Result{}, nil
		}
		
		oldStatus := instance.Status.DeepCopy()
		taskName := instance.Kind + "." + instance.Name
		if result, newStatus, err := s.HandleAgentTaskRound(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), &instance, taskName); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to HandleAgentTaskRound, will retry it, error=%v", err)
			return ctrl.Result{}, err

		} else {
			if newStatus != nil && !reflect.DeepEqual(newStatus, oldStatus) {
				instance.Status = *newStatus
				if err := s.client.Status().Update(ctx, &instance); err != nil {
					// requeue
					logger.Sugar().Errorf("failed to update status, will retry it, error=%v", err)
					return ctrl.Result{}, err
				}
				logger.Sugar().Debugf("succeeded update status, newStatus=%+v", newStatus)
			}

			if result != nil {
				return *result, nil
			}
		}

	default:
		s.logger.Sugar().Errorf("unknown crd type , support kind=%v, detail=%+v", s.crdKind, req)
		// forget this
		return ctrl.Result{}, nil
	}

	return s.plugin.AgentReconcile(s.logger, s.client, ctx, req)
}
