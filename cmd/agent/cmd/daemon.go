// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"github.com/spidernet-io/spiderdoctor/pkg/debug"
	"github.com/spidernet-io/spiderdoctor/pkg/loadRequest"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"time"
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

func testHttp() {
	time.Sleep(30 * time.Second)
	agentServiceUrl := fmt.Sprintf("http://%s.%s.svc.%s:%v/", types.AgentConfig.AgentSerivceIpv4Name, types.AgentConfig.PodNamespace, types.AgentConfig.ClusterDnsDomain, types.AgentConfig.HttpPort)
	qps := 10
	PerRequestTimeoutSecond := 5
	RequestTimeSecond := 2
	rootLogger.Sugar().Infof("send http request to self ipv4 service, url=%s, qps=%v ,PerRequestTimeoutSecond=%v, RequestTimeSecond=%v ", agentServiceUrl, qps, PerRequestTimeoutSecond, RequestTimeSecond)
	r := loadRequest.HttpRequest(agentServiceUrl, qps, PerRequestTimeoutSecond, RequestTimeSecond)
	rootLogger.Sugar().Infof("http result: %v", r)
}

func DaemonMain() {
	rootLogger.Sugar().Infof("config: %+v", types.AgentConfig)

	SetupUtility()

	SetupHttpServer()

	RunMetricsServer(types.AgentConfig.PodName)

	initGrpcServer()

	testHttp()

	rootLogger.Info("hello world")
	time.Sleep(time.Hour)
}
