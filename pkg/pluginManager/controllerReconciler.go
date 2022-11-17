// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type pluginControllerReconciler struct {
	client  client.Client
	plugin  plugintypes.ChainingPlugin
	logger  *zap.Logger
	crdKind string
}

// contorller reconcile
// (1) chedule all task time
// (2) update stauts result
// (3) collect report from agent
func (s *pluginControllerReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	// ------ add crd ------
	switch s.crdKind {
	case KindNameNethttp:
		// ------ add crd ------
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
		if result, newStatus, err := s.UpdateStatus(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), taskName); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to UpdateStatus, will retry it, error=%v", err)
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
		// ------ add crd ------
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
		if result, newStatus, err := s.UpdateStatus(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), taskName); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to UpdateStatus, will retry it, error=%v", err)
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

	}
	// forget this
	return ctrl.Result{}, nil

	// return s.plugin.ControllerReconcile(s.logger, s.client, ctx, req)
}

var _ reconcile.Reconciler = &pluginControllerReconciler{}

func (s *pluginControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).For(s.plugin.GetApiType()).Complete(s)
}
