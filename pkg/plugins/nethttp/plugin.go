// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nethttp

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type pluginNetHttp struct {
}

func (s *pluginNetHttp) GetApiType() client.Object {
	return &crd.Nethttp{}
}

func (s *pluginNetHttp) CheckObjType(obj runtime.Object) bool {
	_, ok := obj.(*crd.Nethttp)
	return ok
}

func (s *pluginNetHttp) ControllerReconcile(client.Client, context.Context, reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func (s *pluginNetHttp) AgentReconcile(client.Client, context.Context, reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func (s *pluginNetHttp) WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	return nil
}

func (s *pluginNetHttp) WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	return nil
}

func (s *pluginNetHttp) WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error {
	return nil
}

func (s *pluginNetHttp) WebhookValidateDelete(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	return nil

}

func init() {
	pluginManager.RegisterPlugin("nethttp", &pluginNetHttp{})
}
