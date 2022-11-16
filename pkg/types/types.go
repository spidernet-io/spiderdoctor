// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package types

type ConfigmapConfig struct {
	EnableIPv4                                     bool `yaml:"enableIPv4"`
	EnableIPv6                                     bool `yaml:"enableIPv6"`
	TaskPollIntervalInSecond                       int  `yaml:"taskPollIntervalInSecond"`
	NethttpDefaultRequestQps                       int  `yaml:"nethttp_defaultRequest_Qps"`
	NethttpDefaultRequestDurationInSecond          int  `yaml:"nethttp_defaultRequest_DurationInSecond"`
	NethttpDefaultRequestPerRequestTimeoutInSecond int  `yaml:"nethttp_defaultRequest_PerRequestTimeoutInSecond"`
	NethttpDefaultFailSuccessRate                  int  `yaml:"nethttp_defaultFail_SuccessRate"`
}

type EnvMapping struct {
	EnvName      string
	DefaultValue string
	P            interface{}
}
