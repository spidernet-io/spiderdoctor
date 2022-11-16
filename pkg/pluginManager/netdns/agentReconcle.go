// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package netdns

import (
	"context"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (s *PluginNetDns) AgentReconcile(logger *zap.Logger, client client.Client, ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// t := &crd.Netdns{
	// 	ObjectMeta: metav1.ObjectMeta{Name: req.Name},
	// }
	// err := client.Get(ctx, req.NamespacedName, t)
	// if err != nil {
	// 	logger.Sugar().Errorf("failed to get netdns %+v", req)
	// 	return ctrl.Result{}, err
	// }
	// logger.Sugar().Infof("get netdns %+v", t)

	return ctrl.Result{}, nil
}
