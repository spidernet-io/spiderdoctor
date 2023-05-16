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

	err := s.RecommendedOptions.ApplyTo(serverConfig)
	if nil != err {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   apiserver.ExtraConfig{},
	}

	return config, nil
}

func NewCommandStartSpiderDoctorServer(stopCh <-chan struct{}) (*cobra.Command, error) {
	options := NewSpiderDoctorServerOptions()

	cmd := &cobra.Command{
		Short: "run a SpiderDoctor api server",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := options.Config()
			if nil != err {
				return err
			}

			server, err := config.Complete().New()
			if nil != err {
				return err
			}
			
			err = server.Run(stopCh)
			if nil != err {
				return err
			}
			return nil
		},
	}

	cmd.Flags()

	return cmd, nil
}
