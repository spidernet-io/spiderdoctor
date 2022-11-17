// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nethttp

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *PluginNetHttp) WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	req, ok := obj.(*crd.Nethttp)
	if !ok {
		s := "failed to get nethttp obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}

	if req.Spec.Target == nil {
		enableIpv4 := types.ControllerConfig.Configmap.EnableIPv4
		enableIpv6 := types.ControllerConfig.Configmap.EnableIPv6
		m := &crd.TargetAgentSepc{
			TestEndpoint:        true,
			TestMultusInterface: true,
			TestClusterIp:       true,
			TestIngress:         false,
			TestIPv6:            &enableIpv4,
			TestIPv4:            &enableIpv6,
			TestNodePort:        true,
		}
		req.Spec.Target = &crd.NethttpTarget{
			TargetAgent: m,
		}
		logger.Sugar().Debugf("set default target for request %v", req.Name)
	}

	if req.Spec.Schedule == nil {
		m := &crd.SchedulePlan{
			StartAfterMinute: 0,
			RoundNumber:      1,
			IntervalMinute:   60,
			TimeoutMinute:    60,
		}
		req.Spec.Schedule = m
		logger.Sugar().Debugf("set default SchedulePlan for request %v", req.Name)
	}

	if req.Spec.Request == nil {
		m := &crd.NethttpRequest{
			DurationInSecond:          types.ControllerConfig.Configmap.NethttpDefaultRequestDurationInSecond,
			QPS:                       types.ControllerConfig.Configmap.NethttpDefaultRequestQps,
			PerRequestTimeoutInSecond: types.ControllerConfig.Configmap.NethttpDefaultRequestPerRequestTimeoutInSecond,
		}
		req.Spec.Request = m
		logger.Sugar().Debugf("set default Request for request %v", req.Name)
	}

	if req.Spec.SuccessCondition == nil {
		m := &crd.NetSuccessCondition{
			SuccessRate:         1,
			MeanAccessDelayInMs: 10000,
		}
		req.Spec.SuccessCondition = m
		logger.Sugar().Debugf("set default SuccessCondition for request %v", req.Name)
	}

	return nil
}

func (s *PluginNetHttp) WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*crd.Nethttp)
	if !ok {
		s := "failed to get nethttp obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Infof("obj: %+v", r)

	// TODO: validate request

	return nil
}

func (s *PluginNetHttp) WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error {
	r, ok := newObj.(*crd.Nethttp)
	if !ok {
		s := "failed to get nethttp obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Infof("obj: %+v", r)

	return nil
}

func (s *PluginNetHttp) WebhookValidateDelete(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	return nil

}
