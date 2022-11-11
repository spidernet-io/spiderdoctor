package grpcManager

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"time"
)

const (
	DefaultDialTimeOut                 = 30 * time.Second
	DefaultClientKeepAliveTimeInterval = 30 * time.Second
	DefaultClientKeepAliveTimeOut      = 10 * time.Second

	LBPolicyFistPick = "pick_first"
	LBPolicyRR       = "round_robin"
)

type GrpcClientManager interface {
	ClientDial(ctx context.Context, serverAddress []string, tlsCaPath string) error
	Close()
}

type grpcClientManager struct {
	logger *zap.Logger
	opts   []grpc.DialOption
	client *grpc.ClientConn
}

func NewGrpcClient(logger *zap.Logger, tlsCaPath string) GrpcClientManager {
	s := &grpcClientManager{
		logger: logger,
	}

	s.opts = append(s.opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(DefaultMaxRecvMsgSize), grpc.MaxCallSendMsgSize(DefaultMaxRecvMsgSize)))

	if len(tlsCaPath) > 0 {
		if creds, err := credentials.NewClientTLSFromFile(tlsCaPath, ""); err != nil {
			s.logger.Sugar().Fatalf("failed to load credentials: %v", err)
		} else {
			s.opts = append(s.opts, grpc.WithTransportCredentials(creds))
		}
	} else {
		s.opts = append(s.opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return s
}

// serverAddress:=[]string{"1.1.1.1:456"}
func (s *grpcClientManager) ClientDial(ctx context.Context, serverAddress []string, tlsCaPath string) error {

	opts := []grpc.DialOption{}

	t := manual.NewBuilderWithScheme("whatever")
	opts = append(opts, grpc.WithResolvers(t))

	serverAddr := t.Scheme() + ":///test.server"
	m := []resolver.Address{}
	for _, address := range serverAddress {
		// https://godoc.org/google.golang.org/grpc/resolver#Address
		m = append(m, resolver.Address{
			Addr: address,
		})
	}
	t.InitialState(resolver.State{Addresses: m})

	// --------
	serviceConfig := map[string]interface{}{}
	serviceConfig["LoadBalancingConfig"] = []map[string]map[string]string{
		{LBPolicyRR: {}},
	}
	serviceConfig["healthCheckConfig"] = map[string]string{
		// An empty string (`""`) typically indicates the overall health of a server should be reported
		"serviceName": "",
	}
	if jsongByte, e := json.Marshal(serviceConfig); e != nil {
		s.logger.Sugar().Fatalf("failed to pase serviceConfig, error=%v", e)
	} else {
		s.logger.Sugar().Debugf("grpc client serviceConfig = %+v \n ", string(jsongByte))
		opts = append(opts, grpc.WithDefaultServiceConfig(string(jsongByte)))
	}

	// --------
	kacp := keepalive.ClientParameters{
		Time:                DefaultClientKeepAliveTimeInterval, // send pings every 10 seconds if there is no activity
		Timeout:             DefaultClientKeepAliveTimeOut,
		PermitWithoutStream: true, // send pings even without active streams
	}
	opts = append(opts, grpc.WithKeepaliveParams(kacp))

	opts = append(opts, s.opts...)

	if c, err := grpc.DialContext(ctx, serverAddr, opts...); err != nil {
		return errors.Errorf("grpc failed to dial")
	} else {
		s.client = c
	}
	return nil
}

func (s *grpcClientManager) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
