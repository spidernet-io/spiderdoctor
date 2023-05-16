package server

import (
	"github.com/spf13/cobra"
	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/pkg/apiserver"
	"github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
	genericapiserver "k8s.io/apiserver/pkg/server"

	genericoptions "k8s.io/apiserver/pkg/server/options"
)

const defaultEtcdPathPrefix = ""

type SpiderDoctorServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions
}

func NewSpiderDoctorServerOptions() *SpiderDoctorServerOptions {
	s := &SpiderDoctorServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			apiserver.Codecs.LegacyCodec(v1beta1.GroupVersion),
		),
	}

	return s
}

func (s *SpiderDoctorServerOptions) Config() (*apiserver.Config, error) {
	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   apiserver.ExtraConfig{},
	}

	return config, nil
}

func NewCommandStartSpiderDoctorServer() (*cobra.Command, error) {

	cmd := &cobra.Command{
		Short: "run a SpiderDoctor api server",
	}

	cmd.Flags()

	return cmd, nil
}
