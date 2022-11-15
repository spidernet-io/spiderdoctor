// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// --------------------

type pluginWebhookhander struct {
	logger *zap.Logger
	plugin plugintypes.ChainingPlugin
	client client.Client
}

var _ webhook.CustomValidator = (*pluginWebhookhander)(nil)

// mutating webhook
func (s *pluginWebhookhander) Default(ctx context.Context, obj runtime.Object) error {
	return s.plugin.WebhookMutating(s.logger.Named("mutatingWebhook"), ctx, obj)
}

func (s *pluginWebhookhander) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	return s.plugin.WebhookValidateCreate(s.logger.Named("validatingCreateWebhook"), ctx, obj)
}

func (s *pluginWebhookhander) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	return s.plugin.WebhookValidateUpdate(s.logger.Named("validatingCreateWebhook"), ctx, oldObj, newObj)

}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (s *pluginWebhookhander) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	return s.plugin.WebhookValidateDelete(s.logger.Named("validatingDeleteWebhook"), ctx, obj)

}

// --------------------

func (s *pluginWebhookhander) SetupWebhook(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(s.plugin.GetApiType()).WithDefaulter(s).WithValidator(s).RecoverPanic().Complete()
}
