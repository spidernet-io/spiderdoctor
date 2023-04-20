// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nethttp

import (
	"context"
	"fmt"
	k8sObjManager "github.com/spidernet-io/spiderdoctor/pkg/k8ObjManager"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/tools"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	networkingv1 "k8s.io/api/networking/v1"
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
		var agentV4Url, agentV6Url *k8sObjManager.ServiceAccessUrl
		var e error

		testIngress := false
		var agentIngress *networkingv1.Ingress
		agentIngress, e = k8sObjManager.GetK8sObjManager().GetIngress(ctx, types.ControllerConfig.Configmap.AgentIngressName, types.ControllerConfig.PodNamespace)
		if e != nil {
			logger.Sugar().Errorf("failed to get ingress , error=%v", e)
		}
		if agentIngress != nil && len(agentIngress.Status.LoadBalancer.Ingress) > 0 {
			testIngress = true
		}

		serviceAccessPortName := "http"
		testLoadBalancer := false
		if types.ControllerConfig.Configmap.EnableIPv4 {
			agentV4Url, e = k8sObjManager.GetK8sObjManager().GetServiceAccessUrl(ctx, types.ControllerConfig.Configmap.AgentSerivceIpv4Name, types.ControllerConfig.PodNamespace, serviceAccessPortName)
			if e != nil {
				logger.Sugar().Errorf("failed to get agent ipv4 service url , error=%v", e)
			}
			if len(agentV4Url.LoadBalancerUrl) > 0 {
				testLoadBalancer = true
			}
		}
		if types.ControllerConfig.Configmap.EnableIPv6 {
			agentV6Url, e = k8sObjManager.GetK8sObjManager().GetServiceAccessUrl(ctx, types.ControllerConfig.Configmap.AgentSerivceIpv4Name, types.ControllerConfig.PodNamespace, serviceAccessPortName)
			if e != nil {
				logger.Sugar().Errorf("failed to get agent ipv6 service url , error=%v", e)
			}
			if len(agentV6Url.LoadBalancerUrl) > 0 {
				testLoadBalancer = true
			}
		}

		enableIpv4 := types.ControllerConfig.Configmap.EnableIPv4
		enableIpv6 := types.ControllerConfig.Configmap.EnableIPv6
		m := &crd.TargetAgentSepc{
			TestEndpoint:        true,
			TestMultusInterface: false,
			TestClusterIp:       true,
			TestNodePort:        true,
			TestLoadBalancer:    testLoadBalancer,
			TestIngress:         testIngress,
			TestIPv6:            &enableIpv6,
			TestIPv4:            &enableIpv4,
		}
		req.Spec.Target = &crd.NethttpTarget{
			TargetAgent: m,
		}
		logger.Sugar().Debugf("set default target for nethttp %v", req.Name)
	}

	if req.Spec.Schedule == nil {
		req.Spec.Schedule = tools.GetDefaultSchedule()
		logger.Sugar().Debugf("set default SchedulePlan for nethttp %v", req.Name)
	}

	if req.Spec.Request == nil {
		m := &crd.NethttpRequest{
			DurationInSecond:      types.ControllerConfig.Configmap.NethttpDefaultRequestDurationInSecond,
			QPS:                   types.ControllerConfig.Configmap.NethttpDefaultRequestQps,
			PerRequestTimeoutInMS: types.ControllerConfig.Configmap.NethttpDefaultRequestPerRequestTimeoutInMS,
		}
		req.Spec.Request = m
		logger.Sugar().Debugf("set default Request for nethttp %v", req.Name)
	}

	if req.Spec.SuccessCondition == nil {
		req.Spec.SuccessCondition = tools.GetDefaultNetSuccessCondition()
		logger.Sugar().Debugf("set default SuccessCondition for nethttp %v", req.Name)
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
		if r.Spec.Request.PerRequestTimeoutInMS > int(r.Spec.Schedule.TimeoutMinute*60*1000) {
			s := fmt.Sprintf("nethttp %v requires PerRequestTimeoutInMS %v ms smaller than Schedule.TimeoutMinute %vm ", r.Name, r.Spec.Request.PerRequestTimeoutInMS, r.Spec.Schedule.TimeoutMinute)
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
		if r.Spec.Target.TargetAgent == nil && r.Spec.Target.TargetUser == nil && r.Spec.Target.TargetPod == nil {
			s := fmt.Sprintf("nethttp %v, no target specified in the spec", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}

		if r.Spec.Target.TargetAgent != nil && r.Spec.Target.TargetUser != nil {
			s := fmt.Sprintf("nethttp %v, forbid to set TargetUser and TargetAgent at same time", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.Target.TargetPod != nil && r.Spec.Target.TargetUser != nil {
			s := fmt.Sprintf("nethttp %v, forbid to set TargetPod and TargetAgent at same time", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.Target.TargetPod != nil && r.Spec.Target.TargetAgent != nil {
			s := fmt.Sprintf("nethttp %v, forbid to set TargetPod and TargetAgent at same time", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}

		if r.Spec.Target.TargetAgent != nil {
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

	return nil
}
