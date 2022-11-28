// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package types

type ConfigmapConfig struct {
	EnableIPv4                                 bool   `yaml:"enableIPv4"`
	EnableIPv6                                 bool   `yaml:"enableIPv6"`
	TaskPollIntervalInSecond                   int    `yaml:"taskPollIntervalInSecond"`
	NethttpDefaultRequestQps                   int    `yaml:"nethttp_defaultRequest_Qps"`
	NethttpDefaultRequestMaxQps                int    `yaml:"nethttp_defaultRequest_MaxQps"`
	NethttpDefaultRequestDurationInSecond      int    `yaml:"nethttp_defaultRequest_DurationInSecond"`
	NethttpDefaultRequestPerRequestTimeoutInMS int    `yaml:"nethttp_defaultRequest_PerRequestTimeoutInMS"`
	MultusPodAnnotationKey                     string `yaml:"multusPodAnnotationKey"`
	CrdMaxHistory                              int    `yaml:"crdMaxHistory"`
}

type EnvMapping struct {
	EnvName      string
	DefaultValue string
	P            interface{}
}
