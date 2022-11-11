package grpc

import (
	"context"
	"github.com/spidernet-io/spiderdoctor/api/v1/agentGrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if len(req.Command) == 0 {
		logger.Error("grpc server ExecRemoteCmd: got empty command \n")
		return nil, status.Error(codes.InvalidArgument, "request command is empty")
	}

	timeout_second := 60
	StdoutMsg, StderrMsg, exitedCode, e := myos.RunCmd(req.Command, nil, "", timeout_second)

	logger.Sugar().Infof("stderrMsg=%v", StderrMsg)
	logger.Sugar().Infof("StdoutMsg=%v", StdoutMsg)
	logger.Sugar().Infof("exitedCode=%v", exitedCode)
	logger.Sugar().Infof("error=%v", e)

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
