package main

import (
	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/cmd/apiserver/server"
	"k8s.io/component-base/cli"

	"k8s.io/klog/v2"
	"os"
)

func main() {
	cmd, err := server.NewCommandStartSpiderDoctorServer()
	if nil != err {
		klog.Errorf("Error creating server: %v", err)
		os.Exit(1)
	}

	code := cli.Run(cmd)
	os.Exit(code)
}
