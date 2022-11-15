// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package types

var ControllerEnvMapping = []EnvMapping{
	{"ENV_ENABLED_METRIC", "false", &ControllerConfig.EnableMetric},
	{"ENV_METRIC_HTTP_PORT", "", &ControllerConfig.MetricPort},
	{"ENV_HTTP_PORT", "80", &ControllerConfig.HttpPort},
	{"ENV_GOPS_LISTEN_PORT", "", &ControllerConfig.GopsPort},
	{"ENV_WEBHOOK_PORT", "", &ControllerConfig.WebhookPort},
	{"ENV_PYROSCOPE_PUSH_SERVER_ADDRESS", "", &ControllerConfig.PyroscopeServerAddress},
	{"ENV_POD_NAME", "", &ControllerConfig.PodName},
	{"ENV_POD_NAMESPACE", "", &ControllerConfig.PodNamespace},
	{"ENV_GOLANG_MAXPROCS", "8", &ControllerConfig.GolangMaxProcs},
	{"ENV_AGGREGATE_AGENT_REPORT", "false", &ControllerConfig.EnableAggregateAgentReport},
	{"ENV_PATH_AGGREGATE_AGENT_REPORT", "/report", &ControllerConfig.DirPathAggregateAgentReport},
	{"ENV_AGENT_GRPC_LISTEN_PORT", "3000", &ControllerConfig.AgentGrpcListenPort},
	{"ENV_PATH_AGENT_POD_REPORT", "/report", &ControllerConfig.AgentPodReportRootDirPath},
	{"ENV_AGENT_DAEMONSET_NAME", "spiderdoctor-agent", &ControllerConfig.SpiderDoctorAgentDaemonsetName},
}

type ControllerConfigStruct struct {
	// ------- from env
	EnableMetric                   bool
	MetricPort                     int32
	HttpPort                       int32
	GopsPort                       int32
	WebhookPort                    int32
	PyroscopeServerAddress         string
	PodName                        string
	PodNamespace                   string
	GolangMaxProcs                 int32
	EnableAggregateAgentReport     bool
	DirPathAggregateAgentReport    string
	AgentGrpcListenPort            int32
	AgentPodReportRootDirPath      string
	SpiderDoctorAgentDaemonsetName string

	// -------- from flags
	ConfigMapPath     string
	TlsCaCertPath     string
	TlsServerCertPath string
	TlsServerKeyPath  string

	// -------- from configmap
	Configmap ConfigmapConfig
}

var ControllerConfig ControllerConfigStruct
