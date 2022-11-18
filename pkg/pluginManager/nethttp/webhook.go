// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nethttp

import (
	"context"
	"fmt"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/tools"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
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

	if req.DeletionTimestamp != nil {
		return nil
	}

	if req.Spec.Target == nil {
		enableIpv4 := types.ControllerConfig.Configmap.EnableIPv4
		enableIpv6 := types.ControllerConfig.Configmap.EnableIPv6
		m := &crd.TargetAgentSepc{
			TestEndpoint:        true,
			TestMultusInterface: true,
			TestClusterIp:       true,
			TestIngress:         false,
			TestIPv6:            &enableIpv6,
			TestIPv4:            &enableIpv4,
			TestNodePort:        true,
		}
		req.Spec.Target = &crd.NethttpTarget{
			TargetAgent: m,
		}
		logger.Sugar().Debugf("set default target for request %v", req.Name)
	}

	if req.Spec.Schedule == nil {
		req.Spec.Schedule = tools.GetDefaultSchedule()
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
		req.Spec.SuccessCondition = tools.GetDefaultNetSuccessCondition()
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
	logger.Sugar().Debugf("Nethttp: %+v", r)

	// validate Schedule
	if true {
		if err := tools.ValidataCrdSchedule(r.Spec.Schedule); err != nil {
			s := fmt.Sprintf("nethttp %v : %v", r.Name, err)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate request
	if true {
		if r.Spec.Request.QPS >= types.ControllerConfig.Configmap.NethttpDefaultRequestMaxQps {
			s := fmt.Sprintf("nethttp %v requires qps %v bigger than maximum %v", r.Name, r.Spec.Request.QPS, types.ControllerConfig.Configmap.NethttpDefaultRequestMaxQps)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.Request.PerRequestTimeoutInSecond > r.Spec.Request.DurationInSecond {
			s := fmt.Sprintf("nethttp %v requires PerRequestTimeoutInSecond %vs smaller than DurationInSecond %vs ", r.Name, r.Spec.Request.PerRequestTimeoutInSecond, r.Spec.Request.DurationInSecond)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.Request.DurationInSecond > int(r.Spec.Schedule.TimeoutMinute*60) {
			s := fmt.Sprintf("nethttp %v requires request.DurationInSecond %vs smaller than Schedule.TimeoutMinute %vm ", r.Name, r.Spec.Request.DurationInSecond, r.Spec.Schedule.TimeoutMinute)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate target
	if true {
		if r.Spec.Target.TargetAgent == nil && r.Spec.Target.TargetUrl == nil {
			s := fmt.Sprintf("nethttp %v, no target specified in the spec", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}

		// validate target
		if r.Spec.Target.TargetAgent.TestIPv4 != nil && *(r.Spec.Target.TargetAgent.TestIPv4) && !types.ControllerConfig.Configmap.EnableIPv4 {
			s := fmt.Sprintf("nethttp %v TestIPv4, but spiderdoctor ipv4 feature is disabled", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.Target.TargetAgent.TestIPv6 != nil && *(r.Spec.Target.TargetAgent.TestIPv6) && !types.ControllerConfig.Configmap.EnableIPv6 {
			s := fmt.Sprintf("nethttp %v TestIPv6, but spiderdoctor ipv6 feature is disabled", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate SuccessCondition
	if true {
		if r.Spec.SuccessCondition.SuccessRate == nil && r.Spec.SuccessCondition.MeanAccessDelayInMs == nil {
			s := fmt.Sprintf("nethttp %v, no SuccessCondition specified in the spec", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.SuccessCondition.SuccessRate != nil && (*(r.Spec.SuccessCondition.SuccessRate) > 1) {
			s := fmt.Sprintf("nethttp %v, SuccessCondition.SuccessRate %v must not be bigger than 1", r.Name, *(r.Spec.SuccessCondition.SuccessRate))
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.SuccessCondition.SuccessRate != nil && (*(r.Spec.SuccessCondition.SuccessRate) < 0) {
			s := fmt.Sprintf("nethttp %v, SuccessCondition.SuccessRate %v must not be smaller than 0 ", r.Name, *(r.Spec.SuccessCondition.SuccessRate))
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	return nil
}

// this will not be called, it is not allowed to modify crd
func (s *PluginNetHttp) WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error {
	r, ok := newObj.(*crd.Nethttp)
	if !ok {
		s := "failed to get nethttp obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Infof("obj: %+v", r.Name)

	if r.DeletionTimestamp == nil {
		return apierrors.NewBadRequest(plugintypes.ApiMsgUnsupportModify)
	}

	return nil
}
