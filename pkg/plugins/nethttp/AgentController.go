// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nethttp

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (s *pluginNetHttp) AgentReconcile(logger *zap.Logger, client client.Client, ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	t := &crd.Nethttp{
		ObjectMeta: metav1.ObjectMeta{Name: req.Name},
	}
	err := client.Get(ctx, req.NamespacedName, t)
	if err != nil {
		logger.Sugar().Errorf("failed to get nethttp %+v", req)
		return reconcile.Result{}, err
	}
	logger.Sugar().Errorf("get nethttp %+v", t)

	return reconcile.Result{}, nil
}
