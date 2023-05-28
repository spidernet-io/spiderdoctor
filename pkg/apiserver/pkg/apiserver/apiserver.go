package apiserver

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/apis/audit/install"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"

	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/pkg/filters"
	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/pkg/registry"
	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/pkg/registry/spidedoctor/pluginreport"
	"github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/system/v1beta1"
)

const DefaultPluginReportPath = "/report"

var (
	Scheme    = runtime.NewScheme()
	Codecs    = serializer.NewCodecFactory(Scheme)
	GroupName = v1beta1.GroupName
)

type ExtraConfig struct {
	DirPathControllerReport string
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

type SpiderDoctorServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return CompletedConfig{&c}
}

func (c completedConfig) New() (*SpiderDoctorServer, error) {
	handlerChainFunc := c.GenericConfig.BuildHandlerChainFunc
	c.GenericConfig.BuildHandlerChainFunc = func(apiHandler http.Handler, c *genericapiserver.Config) http.Handler {
		handler := handlerChainFunc(apiHandler, c)
		handler = filters.WithRequestQuery(handler)
		return handler
	}

	genericServer, err := c.GenericConfig.New("spiderdoctor-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	s := &SpiderDoctorServer{
		GenericAPIServer: genericServer,
	}

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(GroupName, Scheme, metav1.ParameterCodec, Codecs)

	v1beta1storage := map[string]rest.Storage{}
	v1beta1storage["pluginreports"] = registry.RESTInPeace(pluginreport.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter))
	apiGroupInfo.VersionedResourcesStorageMap["v1beta1"] = v1beta1storage

	err = s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo)
	if nil != err {
		return nil, err
	}

	return s, nil
}

func (s *SpiderDoctorServer) Run(stopCh <-chan struct{}) error {
	s.GenericAPIServer.AddPostStartHookOrDie("post-starthook", func(ctx genericapiserver.PostStartHookContext) error {
		return nil
	})

	return s.GenericAPIServer.PrepareRun().Run(stopCh)
}

func init() {
	install.Install(Scheme)
	utilruntime.Must(v1beta1.AddToScheme(Scheme))
	utilruntime.Must(Scheme.SetVersionPriority(v1beta1.SchemeGroupVersion))

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}
