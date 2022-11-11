// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package netdoctorManager

import (
	"github.com/spidernet-io/spiderdoctor/pkg/netdoctorManager/types"
	"go.uber.org/zap"
)

type netdoctorManager struct {
	logger   *zap.Logger
	webhook  *webhookhander
	informer *informerHandler
}

var _ types.NetdoctorManager = (*netdoctorManager)(nil)

func New(logger *zap.Logger) types.NetdoctorManager {

	return &netdoctorManager{
		logger: logger.Named("netdoctorManager"),
	}
}
