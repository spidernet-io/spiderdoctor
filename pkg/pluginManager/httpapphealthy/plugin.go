// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package httpapphealthy

import (
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginHttpAppHealthy struct {
}

var _ types.ChainingPlugin = &PluginHttpAppHealthy{}

func (s *PluginHttpAppHealthy) GetApiType() client.Object {
	return &crd.HttpAppHealthy{}
}
