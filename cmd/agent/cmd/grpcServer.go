// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"github.com/spidernet-io/spiderdoctor/pkg/grpcManager"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"github.com/spidernet-io/spiderdoctor/pkg/utils"
)

const (
	tlsCertPath = "/tmp/cert.crt"
	tlsKeyPath  = "/tmp/key.crt"
	tlsCaPath   = "/tmp/ca.crt"
)

func initGrpcServer() {
	// ---- grpc server
	rootLogger.Info("start grpc server")

	alternateDNS := []string{}
	alternateDNS = append(alternateDNS, types.AgentConfig.PodName)
	// generate self-signed certificates
	if e := utils.NewServerCertKeyForLocalNode(alternateDNS, tlsCertPath, tlsKeyPath, tlsCaPath); e != nil {
		rootLogger.Sugar().Fatalf("failed to generate certiface, error=%v", e)
	}

	t := grpcManager.NewGrpcServer(rootLogger, tlsCertPath, tlsKeyPath)
	listenAddr := fmt.Sprintf(":%d", types.AgentConfig.AgentGrpcListenPort)
	t.Run(listenAddr)
}
