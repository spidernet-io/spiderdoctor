package types

var AgentEnvMapping = []EnvMapping{
	{"ENV_ENABLED_METRIC", "false", &AgentConfig.EnableMetric},
	{"ENV_METRIC_HTTP_PORT", "", &AgentConfig.MetricPort},
	{"ENV_HTTP_PORT", "80", &AgentConfig.HttpPort},
	{"ENV_GOPS_LISTEN_PORT", "", &AgentConfig.GopsPort},
	{"ENV_WEBHOOK_PORT", "", &AgentConfig.WebhookPort},
	{"ENV_PYROSCOPE_PUSH_SERVER_ADDRESS", "", &AgentConfig.PyroscopeServerAddress},
	{"ENV_POD_NAME", "", &AgentConfig.PodName},
	{"ENV_POD_NAMESPACE", "", &AgentConfig.PodNamespace},
	{"ENV_GOLANG_MAXPROCS", "8", &AgentConfig.GolangMaxProcs},
	{"ENV_GC_REPORT_FOR_DELETED_CRD", "true", &AgentConfig.GcReportForDeletdCRD},
	{"ENV_AGENT_GRPC_LISTEN_PORT", "3000", &AgentConfig.AgentGrpcListenPort},
	{"ENV_PATH_AGENT_POD_REPORT", "/report", &AgentConfig.ReportRootDirPath},
}

type AgentConfigStruct struct {
	// ------- from env
	EnableMetric           bool
	MetricPort             int32
	HttpPort               int32
	GopsPort               int32
	WebhookPort            int32
	PyroscopeServerAddress string
	PodName                string
	PodNamespace           string
	GolangMaxProcs         int32
	GcReportForDeletdCRD   bool
	AgentGrpcListenPort    int32
	ReportRootDirPath      string

	// ------- from flags
	ConfigMapPath     string
	TlsCaCertPath     string
	TlsServerCertPath string
	TlsServerKeyPath  string

	// from configmap
	Configmap ConfigmapConfig
}

var AgentConfig AgentConfigStruct
