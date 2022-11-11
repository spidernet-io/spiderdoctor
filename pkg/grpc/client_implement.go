package grpc

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spidernet-io/spiderdoctor/api/v1/agentGrpc"
)

func (s *grpcClientManager) SendRequestForExecRequest(ctx context.Context, request *agentGrpc.ExecRequestMsg) (*agentGrpc.ExecResponseMsg, error) {
	if s.client == nil {
		return nil, errors.Errorf("please dial first")
	}
	c := agentGrpc.NewCmdServiceClient(s.client)

	if r, err := c.ExecRemoteCmd(ctx, request); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}
