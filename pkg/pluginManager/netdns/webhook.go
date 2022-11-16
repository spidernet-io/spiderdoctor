// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package netdns

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *PluginNetDns) WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*crd.Netdns)
	if !ok {
		s := "failed to get netnds obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Infof("obj: %+v", r)

	// TODO: mutating defaul value

	return nil
}

func (s *PluginNetDns) WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*crd.Netdns)
	if !ok {
		s := "failed to get netdns obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Infof("obj: %+v", r)

	// TODO: validate request

	return nil
}

func (s *PluginNetDns) WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error {
	r, ok := newObj.(*crd.Netdns)
	if !ok {
		s := "failed to get netdns obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Infof("obj: %+v", r)

	return nil
}

func (s *PluginNetDns) WebhookValidateDelete(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	return nil

}
