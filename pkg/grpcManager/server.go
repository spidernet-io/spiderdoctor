package grpcManager

import (
	"crypto/tls"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"math"
	"net"
	"time"
)

type GrpcServerManager interface {
	Run(listenAddress string)
	UpdateHealthStatus(status healthpb.HealthCheckResponse_ServingStatus)
}

type grpcServer struct {
	server             *grpc.Server
	healthcheckService *health.Server
	logger             *zap.Logger
}

const (
	DefaultServerKeepAliveTimeInterval      = 60 * time.Second
	DefaultServerKeepAliveTimeOut           = 10 * time.Second
	DefaultServerKeepAlivedAlowPeerInterval = 5 * time.Second
	// 100 MByte
	DefaultMaxRecvMsgSize = 1024 * 1024 * 100
	DefaultMaxSendMsgSize = math.MaxInt32
)

func NewGrpcServer(logger *zap.Logger, tlsCaPath, tlsCertPath, tlskeyPath string) GrpcServerManager {
	m := &grpcServer{}
	opts := []grpc.ServerOption{}

	m.logger = logger

	// ----- tls
	cert, err := tls.LoadX509KeyPair(tlsCertPath, tlskeyPath)
	if err != nil {
		logger.Sugar().Fatalf("failed to load tls: %s", err)
	}
	opts = append(opts, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))

	// https://godoc.org/google.golang.org/grpc/keepalive#EnforcementPolicy
	// Enforcement policy is a special setting on server side to protect server from malicious or misbehaving clients
	// for case : (1)Client sends too frequent pings (2)Client sends pings when there's no stream and this is disallowed by server config
	kaep := keepalive.EnforcementPolicy{
		MinTime:             DefaultServerKeepAlivedAlowPeerInterval,
		PermitWithoutStream: true,
	}
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(kaep))

	// https://godoc.org/google.golang.org/grpc/keepalive#ServerParameters
	kasp := keepalive.ServerParameters{
		Time:    DefaultServerKeepAliveTimeInterval,
		Timeout: DefaultServerKeepAliveTimeOut,
	}
	// https://godoc.org/google.golang.org/grpc#KeepaliveParams
	opts = append(opts, grpc.KeepaliveParams(kasp))

	opts = append(opts, grpc.MaxRecvMsgSize(DefaultMaxRecvMsgSize))
	opts = append(opts, grpc.MaxSendMsgSize(DefaultMaxSendMsgSize))

	m.server = grpc.NewServer(opts...)
	if m.server == nil {
		logger.Fatal("failed to New Grpc Server ")
	}
	m.healthcheckService = health.NewServer()
	healthpb.RegisterHealthServer(m.server, m.healthcheckService)
	reflection.Register(m.server)

	m.registerService()

	return m
}

// address: "127.0.0.1:5000" or ":5000"
func (t *grpcServer) Run(listenAddress string) {

	d, err := net.Listen("tcp", listenAddress)
	if err != nil {
		t.logger.Sugar().Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := t.server.Serve(d); err != nil {
			t.logger.Sugar().Fatalf("failed to run Grpc Server, reason=%v", err)
		}
	}()

}

// type HealthCheckResponse_ServingStatus int32
// const (
//
//	HealthCheckResponse_UNKNOWN         HealthCheckResponse_ServingStatus = 0
//	HealthCheckResponse_SERVING         HealthCheckResponse_ServingStatus = 1
//	HealthCheckResponse_NOT_SERVING     HealthCheckResponse_ServingStatus = 2
//	HealthCheckResponse_SERVICE_UNKNOWN HealthCheckResponse_ServingStatus = 3 // Used only by the Watch method.
//
// )
func (t *grpcServer) UpdateHealthStatus(status healthpb.HealthCheckResponse_ServingStatus) {
	if t.healthcheckService == nil {
		return
	}

	t.healthcheckService.SetServingStatus("", status)
	t.logger.Sugar().Infof("grpc server update health status to %v", status)
	return
}
