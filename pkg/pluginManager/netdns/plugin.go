// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package netdns

import (
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginNetDns struct {
}

var _ types.ChainingPlugin = &PluginNetDns{}

func (s *PluginNetDns) GetApiType() client.Object {
	return &crd.Netdns{}
}

func (s *PluginNetDns) GetKindName() string {
	return "Netdns"
}
