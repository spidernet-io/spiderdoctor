// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	"fmt"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"time"
)

// --------------------

type webhookhander struct {
	logger *zap.Logger
	plugin plugintypes.ChainingPlugin
}

var _ webhook.CustomValidator = (*webhookhander)(nil)

// mutating webhook
func (s *webhookhander) Default(ctx context.Context, obj runtime.Object) error {
	return s.plugin.WebhookMutating(s.logger.Named("mutatingWebhook"), ctx, obj)
}

func (s *webhookhander) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	return s.plugin.WebhookValidateCreate(s.logger.Named("validatingCreateWebhook"), ctx, obj)
}

func (s *webhookhander) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	return s.plugin.WebhookValidateUpdate(s.logger.Named("validatingCreateWebhook"), ctx, oldObj, newObj)

}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (s *webhookhander) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	return s.plugin.WebhookValidateDelete(s.logger.Named("validatingDeleteWebhook"), ctx, obj)

}

// --------------------

func (s *pluginManager) runWebhook(webhookPort int, webhookTlsDir string) {

	logger := s.logger
	scheme := runtime.NewScheme()
	if e := plugin.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add scheme for plugin, reason=%v", name, e)
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     "0",
		HealthProbeBindAddress: "0",
		// webhook port
		Port: webhookPort,
		// directory that contains the webhook server key and certificate, The server key and certificate must be named tls.key and tls.crt
		CertDir: webhookTlsDir,
	})
	if err != nil {
		logger.Sugar().Fatalf("failed to NewManager, reason=%v", err)
	}

	for name, plugin := range s.chainingPlugins {
		logger.Sugar().Infof("setup webhook for plugin %v on port %v, with tls under %v", name, webhookPort, webhookTlsDir)

		eqw := &webhookhander{
			logger: logger,
			plugin: plugin,
		}
		e := ctrl.NewWebhookManagedBy(mgr).For(plugin.GetApiType()).WithDefaulter(eqw).WithValidator(eqw).RecoverPanic().Complete()
		if e != nil {
			logger.Sugar().Fatalf("failed to NewWebhookManagedBy, reason=%v", e)
		}
	}
	go func() {
		s := "webhook down"
		// mgr.Start()
		if err := mgr.Start(context.Background()); err != nil {
			s += fmt.Sprintf(", reason=%v", err)
		}
		logger.Error(s)
		time.Sleep(time.Second)
	}()

}
