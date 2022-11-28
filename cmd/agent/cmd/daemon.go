// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spidernet-io/spiderdoctor/pkg/debug"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
)

func SetupUtility() {
	// run gops
	d := debug.New(rootLogger)
	if types.AgentConfig.GopsPort != 0 {
		d.RunGops(int(types.AgentConfig.GopsPort))
	}

	if types.AgentConfig.PyroscopeServerAddress != "" {
		d.RunPyroscope(types.AgentConfig.PyroscopeServerAddress, types.AgentConfig.PodName)
	}
}

func DaemonMain() {
	rootLogger.Sugar().Infof("config: %+v", types.AgentConfig)

	if types.AgentConfig.AppMode {
		// app mode, just used to debug
		rootLogger.Info("run in app mode")
		SetupHttpServer()
		initGrpcServer()
		// sleep forever
		select {}
	}

	SetupUtility()

	SetupHttpServer()

	RunMetricsServer(types.AgentConfig.PodName)

	initGrpcServer()

	s := pluginManager.InitPluginManager(rootLogger.Named("agentContorller"))
	s.RunAgentController()

	rootLogger.Info("finish initialization")
	// sleep forever
	select {}
}
