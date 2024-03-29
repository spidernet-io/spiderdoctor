// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/httpapphealthy"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/netdns"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/netreachhealthy"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
)

var pluginLock = &lock.Mutex{}

type pluginManager struct {
	chainingPlugins map[string]plugintypes.ChainingPlugin
	logger          *zap.Logger
}
type PluginManager interface {
	RunAgentController()
	RunControllerController(healthPort int, webhookPort int, webhookTlsDir string)
}

var globalPluginManager *pluginManager

// -------------------------

func InitPluginManager(logger *zap.Logger) PluginManager {
	pluginLock.Lock()
	defer pluginLock.Unlock()

	globalPluginManager.logger = logger

	return globalPluginManager
}

const (
	// ------ add crd ------
	KindNameHttpAppHealthy  = "HttpAppHealthy"
	KindNameNetReachHealthy = "NetReachHealthy"
	KindNameNetdns          = "Netdns"
)

func init() {
	globalPluginManager = &pluginManager{
		chainingPlugins: map[string]plugintypes.ChainingPlugin{},
	}

	// ------ add crd ------
	globalPluginManager.chainingPlugins[KindNameHttpAppHealthy] = &httpapphealthy.PluginHttpAppHealthy{}
	globalPluginManager.chainingPlugins[KindNameNetReachHealthy] = &netreachhealthy.PluginNetReachHealthy{}
	globalPluginManager.chainingPlugins[KindNameNetdns] = &netdns.PluginNetDns{}

}
