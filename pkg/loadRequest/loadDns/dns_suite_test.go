// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0
package loadDns_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	config "github.com/spidernet-io/spiderdoctor/pkg/types"
)

func TestIppoolCR(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "load request Suite")
}

var _ = BeforeSuite(func() {
	// nothing to do
	config.AgentConfig.Configmap.NetdnsDefaultConcurrency = 50
})
