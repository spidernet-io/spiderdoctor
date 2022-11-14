// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nethttp

import (
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginNetHttp struct {
}

var _ types.ChainingPlugin = &PluginNetHttp{}

func (s *PluginNetHttp) GetApiType() client.Object {
	return &crd.Nethttp{}
}

func (s *PluginNetHttp) AddToScheme(t *runtime.Scheme) error {
	return crd.AddToScheme(t)
}
