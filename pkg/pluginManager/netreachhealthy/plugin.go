// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package netreachhealthy

import (
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginNetReachHealthy struct {
}

var _ types.ChainingPlugin = &PluginNetReachHealthy{}

func (s *PluginNetReachHealthy) GetApiType() client.Object {
	return &crd.NetReachHealthy{}
}
