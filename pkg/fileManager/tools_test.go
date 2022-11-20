// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0
package fileManager_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/spiderdoctor/pkg/fileManager"
)

var _ = Describe("test ippool CR", Label("ippoolCR"), func() {

	It("test basic", func() {
		filePath := "/tmp/_loggertest/a.txt"

		wr := fileManager.NewFileWriter(filePath)
		GinkgoWriter.Printf("succeed to new write for %v", filePath)
		defer wr.Close()

		data := []byte("test data\n dsf\n")
		_, t := wr.Write(data)
		Expect(t).NotTo(HaveOccurred())

	})

})
