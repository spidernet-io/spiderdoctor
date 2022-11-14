// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nethttp

import (
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func init() {
	pluginManager.RegisterPlugin("nethttp", &pluginNetHttp{})
}
