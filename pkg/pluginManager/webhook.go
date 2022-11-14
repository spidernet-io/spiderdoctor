// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"time"
)

// --------------------

type webhookhander struct {
	logger *zap.Logger
}

var _ webhook.CustomValidator = (*webhookhander)(nil)

// mutating webhook
func (s *webhookhander) Default(ctx context.Context, obj runtime.Object) error {
	logger := s.logger.Named("mutatingWebhook")
	for name, p := range globalPluginManager.chainingPlugins {
		if p.CheckObjType(obj) {
			return p.WebhookMutating(logger.Named(name), ctx, obj)
		}
	}
	msg := fmt.Sprintf("failed to find plugin to handle obj type %v", reflect.TypeOf(obj))
	logger.Error(msg)
	return fmt.Errorf(msg)

}

func (s *webhookhander) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	logger := s.logger.Named("validatingCreateWebhook")
	for name, p := range globalPluginManager.chainingPlugins {
		if p.CheckObjType(obj) {
			return p.WebhookValidateCreate(logger.Named(name), ctx, obj)
		}
	}
	msg := fmt.Sprintf("failed to find plugin to handle obj type %v", reflect.TypeOf(obj))
	logger.Error(msg)
	return fmt.Errorf(msg)
}

func (s *webhookhander) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	logger := s.logger.Named("validatingUpdatingWebhook")
	for name, p := range globalPluginManager.chainingPlugins {
		if p.CheckObjType(newObj) {
			return p.WebhookValidateUpdate(logger.Named(name), ctx, oldObj, newObj)
		}
	}
	msg := fmt.Sprintf("failed to find plugin to handle obj type %v", reflect.TypeOf(newObj))
	logger.Error(msg)
	return fmt.Errorf(msg)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (s *webhookhander) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	logger := s.logger.Named("validatingDeletingWebhook")
	for name, p := range globalPluginManager.chainingPlugins {
		if p.CheckObjType(obj) {
			return p.WebhookValidateDelete(logger.Named(name), ctx, obj)
		}
	}
	msg := fmt.Sprintf("failed to find plugin to handle obj type %v", reflect.TypeOf(obj))
	logger.Error(msg)
	return fmt.Errorf(msg)
}

// --------------------

func (s *pluginManager) runWebhook(webhookPort int, webhookTlsDir string) {

	logger := s.logger
	r := &webhookhander{
		logger: logger,
	}

	logger.Sugar().Infof("setup webhook on port %v, with tls under %v", webhookPort, webhookTlsDir)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
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

	t := ctrl.NewWebhookManagedBy(mgr)
	for name, p := range s.chainingPlugins {
		logger.Sugar().Infof("run webbook for plugin %v", name)
		t = t.For(p.GetApiType())
	}
	e := t.WithDefaulter(r).WithValidator(r).RecoverPanic().Complete()
	if e != nil {
		logger.Sugar().Fatalf("failed to NewWebhookManagedBy, reason=%v", e)
	}

	go func() {
		s := "webhook down"

		// mgr.Start()
		if err := mgr.GetWebhookServer().Start(context.Background()); err != nil {
			s += fmt.Sprintf(", reason=%v", err)
		}
		logger.Error(s)
		time.Sleep(time.Second)
	}()

}
