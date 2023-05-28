package main

import (
	"os"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"
	"k8s.io/klog/v2"

	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/cmd/apiserver/server"
)

func main() {
	stopCh := genericapiserver.SetupSignalHandler()
	cmd, err := server.NewCommandStartSpiderDoctorServer(stopCh)
	if nil != err {
		klog.Errorf("Error creating server: %v", err)
		os.Exit(1)
	}

	code := cli.Run(cmd)
	os.Exit(code)
}
