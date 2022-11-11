package grpcManager

import (
	"context"
	"github.com/spidernet-io/spiderdoctor/api/v1/agentGrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"github.com/spidernet-io/spiderdoctor/pkg/utils"
)

// ------ implement
type myGrpcServer struct {
	agentGrpc.UnimplementedCmdServiceServer
	logger *zap.Logger
}

func (s *myGrpcServer) ExecRemoteCmd(ctx context.Context, req *agentGrpc.ExecRequestMsg) (*agentGrpc.ExecResponseMsg, error) {

	logger := s.logger.With(
		zap.String("SubnetName", req.Command),
	)
	logger.Sugar().Infof("request: %+v", req)

	if len(req.Command) == 0 {
		logger.Error("grpc server ExecRemoteCmd: got empty command \n")
		return nil, status.Error(codes.InvalidArgument, "request command is empty")
	}
	if req.Timeoutsecond == 0 {
		logger.Error("grpc server ExecRemoteCmd: got empty timeout \n")
		return nil, status.Error(codes.InvalidArgument, "request command is empty")
	}

	clientctx, cancel := context.WithTimeout(context.Background(), time.Duration(req.Timeoutsecond)*time.Second)
	defer cancel()
	go func() {
		select {
		case <-clientctx.Done():
		case <-ctx.Done():
			cancel()
		}
	}()

	StdoutMsg, StderrMsg, exitedCode, e := utils.RunFrondendCmd(clientctx, req.Command, nil, "")

	logger.Sugar().Infof("stderrMsg = %v", StderrMsg)
	logger.Sugar().Infof("StdoutMsg = %v", StdoutMsg)
	logger.Sugar().Infof("exitedCode = %v", exitedCode)
	logger.Sugar().Infof("error = %v", e)

	b := agentGrpc.ExecResponseMsg{
		Stdmsg: StdoutMsg,
		Stderr: StderrMsg,
		Code:   int32(exitedCode),
	}
	return &b, nil
}

// ------------
func (t *grpcServer) registerService() {
	agentGrpc.RegisterCmdServiceServer(t.server, &myGrpcServer{
		logger: t.logger,
	})
}
