package pluginreport

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend/factory"
	"k8s.io/klog/v2"

	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/pkg/registry"
	"github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
)

func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*registry.REST, error) {
	strategy := NewStrategy(scheme)

	restOptions, err := optsGetter.GetRESTOptions(v1beta1.Resource("pluginreports"))
	if nil != err {
		return nil, err
	}

	dryRunnableStorage, destroyFunc := NewStorage(restOptions)
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &v1beta1.PluginReport{} },
		NewListFunc:              func() runtime.Object { return &v1beta1.PluginReportList{} },
		DefaultQualifiedResource: v1beta1.Resource("pluginreports"),
		PredicateFunc:            MatchPluginReport,
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
		Storage:                  dryRunnableStorage,
		DestroyFunc:              destroyFunc,
		TableConvertor:           rest.NewDefaultTableConvertor(v1beta1.Resource("pluginreports")),
	}

	return &registry.REST{Store: store}, nil
}

func NewStorage(restOptions generic.RESTOptions) (genericregistry.DryRunnableStorage, factory.DestroyFunc) {

	dryRunnableStorage := genericregistry.DryRunnableStorage{
		Storage: nil,
		Codec:   restOptions.StorageConfig.Codec,
	}

	return dryRunnableStorage, func() {}
}

var _ storage.Interface = &pluginReportStorage{}

type pluginReportStorage struct {
	resourceName string
}

func (p pluginReportStorage) Versioner() storage.Versioner {
	return storage.APIObjectVersioner{}
}

func (p pluginReportStorage) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {
	return fmt.Errorf("create API not implement")
}

func (p pluginReportStorage) Delete(ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions, validateDeletion storage.ValidateObjectFunc, cachedExistingObject runtime.Object) error {
	return fmt.Errorf("delete API not implement")
}

func (p pluginReportStorage) Watch(ctx context.Context, key string, opts storage.ListOptions) (watch.Interface, error) {
	return nil, fmt.Errorf("watch API not implement")

}

func (p pluginReportStorage) Get(ctx context.Context, key string, opts storage.GetOptions, objPtr runtime.Object) error {
	klog.Infof("Get called with key: %v on resource %v\n", key, p.resourceName)

	pluginReport := objPtr.(*v1beta1.PluginReport)
	pluginReport.TypeMeta = metav1.TypeMeta{
		Kind:       "PluginReport",
		APIVersion: v1beta1.GroupVersion.String(),
	}
	pluginReport.ObjectMeta = metav1.ObjectMeta{
		Name: "test-wk",
	}

	return nil
}

func (p pluginReportStorage) GetList(ctx context.Context, key string, opts storage.ListOptions, listObj runtime.Object) error {
	return fmt.Errorf("GetList API not implement")
}

func (p pluginReportStorage) GuaranteedUpdate(ctx context.Context, key string, destination runtime.Object, ignoreNotFound bool, preconditions *storage.Preconditions, tryUpdate storage.UpdateFunc, cachedExistingObject runtime.Object) error {
	return fmt.Errorf("GuaranteedUpdate API not implement")
}

func (p pluginReportStorage) Count(key string) (int64, error) {
	return 0, fmt.Errorf("Count not supported for key: %s", key)
}
