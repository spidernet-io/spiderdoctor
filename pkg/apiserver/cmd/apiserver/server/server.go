package server

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"

	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/pkg/apiserver"
	"github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/system/v1beta1"
)

const defaultEtcdPathPrefix = ""

type SpiderDoctorServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions
}

func NewSpiderDoctorServerOptions() *SpiderDoctorServerOptions {
	s := &SpiderDoctorServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			apiserver.Codecs.LegacyCodec(v1beta1.SchemeGroupVersion),
		),
	}
	s.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = runtime.NewMultiGroupVersioner(v1beta1.SchemeGroupVersion, schema.GroupKind{Group: v1beta1.GroupName})

	return s
}

func (s *SpiderDoctorServerOptions) Config() (*apiserver.Config, error) {
	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)

	err := s.RecommendedOptions.ApplyTo(serverConfig)
	if nil != err {
		return nil, err
	}

	pluginReportDir := apiserver.DefaultPluginReportPath
	env, ok := os.LookupEnv("ENV_CONTROLLER_REPORT_STORAGE_PATH")
	if ok {
		pluginReportDir = env
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig: apiserver.ExtraConfig{
			DirPathControllerReport: pluginReportDir,
		},
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

	flags := cmd.Flags()
	options.RecommendedOptions.AddFlags(flags)
	utilfeature.DefaultMutableFeatureGate.AddFlag(flags)

	return cmd, nil
}
