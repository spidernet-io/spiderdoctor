// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"context"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ChainingPlugin interface {
	GetApiType() client.Object
	GetKindName() string
	ControllerReconcile(*zap.Logger, client.Client, context.Context, reconcile.Request) (reconcile.Result, error)
	AgentReconcile(*zap.Logger, client.Client, context.Context, reconcile.Request) (reconcile.Result, error)

	WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
	WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
	WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error
	WebhookValidateDelete(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
}
