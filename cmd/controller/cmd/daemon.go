// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"github.com/spidernet-io/spiderdoctor/pkg/debug"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.opentelemetry.io/otel/attribute"
	"path/filepath"
	"time"
)

func SetupUtility() {

	// run gops
	d := debug.New(rootLogger)
	if types.ControllerConfig.GopsPort != 0 {
		d.RunGops(int(types.ControllerConfig.GopsPort))
	}

	if types.ControllerConfig.PyroscopeServerAddress != "" {
		d.RunPyroscope(types.ControllerConfig.PyroscopeServerAddress, types.ControllerConfig.PodName)
	}
}

func DaemonMain() {

	rootLogger.Sugar().Infof("config: %+v", types.ControllerConfig)

	SetupUtility()

	SetupHttpServer()

	// ------

	RunMetricsServer(types.ControllerConfig.PodName)
	MetricGaugeEndpoint.Add(context.Background(), 100)
	MetricGaugeEndpoint.Add(context.Background(), -10)
	MetricGaugeEndpoint.Add(context.Background(), 5)

	attrs := []attribute.KeyValue{
		attribute.Key("pod1").String("value1"),
	}
	MetricCounterRequest.Add(context.Background(), 10, attrs...)
	attrs = []attribute.KeyValue{
		attribute.Key("pod2").String("value1"),
	}
	MetricCounterRequest.Add(context.Background(), 5, attrs...)

	MetricHistogramDuration.Record(context.Background(), 10)
	MetricHistogramDuration.Record(context.Background(), 20)

	// ----------
	s := pluginManager.InitPluginManager(rootLogger.Named("pluginsManager"))
	s.RunControllerController(int(types.ControllerConfig.WebhookPort), filepath.Dir(types.ControllerConfig.TlsServerCertPath))

	// ------------
	rootLogger.Info("hello world")
	time.Sleep(time.Hour)
}
