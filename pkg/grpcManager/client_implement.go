// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package grpcManager

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spidernet-io/spiderdoctor/api/v1/agentGrpc"
)

func (s *grpcClientManager) SendRequestForExecRequest(ctx context.Context, serverAddress []string, request *agentGrpc.ExecRequestMsg) (*agentGrpc.ExecResponseMsg, error) {

	if e := s.clientDial(ctx, serverAddress); e != nil {
		return nil, errors.Errorf("failed to dial, error=%v", e)
	}
	defer s.client.Close()

	c := agentGrpc.NewCmdServiceClient(s.client)

	if r, err := c.ExecRemoteCmd(ctx, request); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}
