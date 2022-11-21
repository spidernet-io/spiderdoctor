// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package grpcManager

import (
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spidernet-io/spiderdoctor/api/v1/agentGrpc"
	"os"
	"strings"
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

func (s *grpcClientManager) GetFileList(ctx context.Context, serverAddress, directory string) ([]string, error) {
	// get agent files list
	request := &agentGrpc.ExecRequestMsg{
		Timeoutsecond: 10,
		Command:       fmt.Sprintf("ls %v", directory),
	}
	response, e := s.SendRequestForExecRequest(ctx, []string{serverAddress}, request)
	if e != nil {
		return nil, fmt.Errorf("failed to get file list under directory %v of %v, error=%v", directory, serverAddress, e)
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("failed to get file list under directory %v of %v, exit code=%v, stderr=%v", directory, serverAddress, response.Code, response.Stderr)
	}

	return strings.Fields(response.Stdmsg), nil
}

func (s *grpcClientManager) SaveRemoteFileToLocal(ctx context.Context, serverAddress, remoteFilePath, localFilePath string) error {

	// get agent files list
	request := &agentGrpc.ExecRequestMsg{
		Timeoutsecond: 10,
		Command:       fmt.Sprintf("cat %v", remoteFilePath),
	}
	response, e := s.SendRequestForExecRequest(ctx, []string{serverAddress}, request)
	if e != nil {
		return fmt.Errorf("failed to get remote file %v of %v, error=%v", remoteFilePath, serverAddress, e)
	}
	if response.Code != 0 {
		return fmt.Errorf("failed to get remote file %v of %v, exit code=%v, stderr=%v", remoteFilePath, serverAddress, response.Code, response.Stderr)
	}

	if len(response.Stdmsg) == 0 {
		return fmt.Errorf("got empty remote file %v of %v ", remoteFilePath, serverAddress)
	}

	f, e := os.OpenFile(localFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if e != nil {
		return fmt.Errorf("open file %v failed, error=%v", localFilePath, e)
	}
	defer f.Close()
	// _, e = f.Write([]byte(response.Stdmsg))
	// if e != nil {
	// 	return fmt.Errorf("failed to write file %v, error=%v", localFilePath, e)
	// }

	writer := bufio.NewWriter(f)
	_, e = writer.WriteString(response.Stdmsg)
	if e != nil {
		return fmt.Errorf("failed to write file %v, error=%v", localFilePath, e)
	}
	if e := writer.Flush(); e != nil {
		return fmt.Errorf("failed to flush file %v, error=%v", localFilePath, e)
	}

	return nil
}
