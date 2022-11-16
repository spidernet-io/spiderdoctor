// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nethttp

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *PluginNetHttp) WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*crd.Nethttp)
	if !ok {
		s := "failed to get nethttp obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Infof("obj: %+v", r)

	// TODO: mutating defaul value

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
