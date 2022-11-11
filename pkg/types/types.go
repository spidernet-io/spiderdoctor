// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package types

type ConfigmapConfig struct {
	EnableIPv4 bool `yaml:"enableIPv4"`
	EnableIPv6 bool `yaml:"enableIPv6"`
}
