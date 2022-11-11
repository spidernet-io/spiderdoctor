// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"time"
)

func DaemonMain() {

	rootLogger.Sugar().Infof("config: %+v", globalConfig)

	rootLogger.Info("hello world")
	time.Sleep(time.Hour)
}
