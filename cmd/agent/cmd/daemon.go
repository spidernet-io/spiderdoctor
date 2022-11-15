// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spidernet-io/spiderdoctor/pkg/debug"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager"
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

/*
func testHttp() {
	time.Sleep(30 * time.Second)
	agentServiceUrl := fmt.Sprintf("http://%s.%s.svc.%s:%v/", types.AgentConfig.AgentSerivceIpv4Name, types.AgentConfig.PodNamespace, types.AgentConfig.ClusterDnsDomain, types.AgentConfig.HttpPort)
	qps := 10
	PerRequestTimeoutSecond := 5
	RequestTimeSecond := 2
	rootLogger.Sugar().Infof("send http request to self ipv4 service, url=%s, qps=%v ,PerRequestTimeoutSecond=%v, RequestTimeSecond=%v ", agentServiceUrl, qps, PerRequestTimeoutSecond, RequestTimeSecond)
	r := loadRequest.HttpRequest(agentServiceUrl, qps, PerRequestTimeoutSecond, RequestTimeSecond)
	if r.Requests == 0 {
		rootLogger.Sugar().Error("failed to send request")
	} else {
		if r.Success == 0 {
			rootLogger.Error("network failed to reach")
		} else if r.Success != 0 {
			rootLogger.Sugar().Warnf("partial request failed")
		} else {
			rootLogger.Info("all request succeed")
		}
	}
	rootLogger.Sugar().Infof("http result: %v", r)
	rootLogger.Sugar().Infof("http Requests Success rate: %v", r.Success)
	rootLogger.Sugar().Infof("http Requests total: %v", r.Requests)
	rootLogger.Sugar().Infof("http Latencies.Mean: %v", r.Latencies.Mean)
	rootLogger.Sugar().Infof("http Latencies.Max: %v", r.Latencies.Max)
	rootLogger.Sugar().Infof("http Latencies.Min: %v", r.Latencies.Min)
	rootLogger.Sugar().Infof("http Duration : %v", r.Duration)
	rootLogger.Sugar().Infof("http sent requests per second : %v", r.Rate)

}
*/

func DaemonMain() {
	rootLogger.Sugar().Infof("config: %+v", types.AgentConfig)

	SetupUtility()

	SetupHttpServer()

	RunMetricsServer(types.AgentConfig.PodName)

	initGrpcServer()

	// testHttp()

	s := pluginManager.InitPluginManager(rootLogger.Named("agentContorller"))
	s.RunAgentController()

	rootLogger.Info("hello world")
	time.Sleep(time.Hour)
}
