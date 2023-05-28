package pluginreport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend/factory"
	"k8s.io/klog/v2"

	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/pkg/printers"
	"github.com/spidernet-io/spiderdoctor/pkg/apiserver/pkg/registry"
	"github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/system/v1beta1"
)

func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*registry.REST, error) {
	strategy := NewStrategy(scheme)

	restOptions, err := optsGetter.GetRESTOptions(v1beta1.Resource("pluginreports"))
	if nil != err {
		return nil, err
	}

	dryRunnableStorage, destroyFunc := NewStorage(restOptions)
	store := &genericregistry.Store{
		NewFunc:     func() runtime.Object { return &v1beta1.PluginReport{} },
		NewListFunc: func() runtime.Object { return &v1beta1.PluginReportList{} },
		KeyRootFunc: func(ctx context.Context) string {
			return restOptions.ResourcePrefix
		},
		KeyFunc: func(ctx context.Context, name string) (string, error) {
			return genericregistry.NoNamespaceKeyFunc(ctx, restOptions.ResourcePrefix, name)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*v1beta1.PluginReport).Name, nil
		},
		DefaultQualifiedResource: v1beta1.Resource("pluginreports"),
		PredicateFunc:            MatchPluginReport,

		CreateStrategy:          strategy,
		UpdateStrategy:          strategy,
		DeleteStrategy:          strategy,
		EnableGarbageCollection: true,

		Storage:     dryRunnableStorage,
		DestroyFunc: destroyFunc,

		TableConvertor: printers.NewTableGenerator(v1beta1.Resource("pluginreports")),
	}

	return &registry.REST{Store: store}, nil
}

func NewStorage(restOptions generic.RESTOptions) (genericregistry.DryRunnableStorage, factory.DestroyFunc) {

	dryRunnableStorage := genericregistry.DryRunnableStorage{
		Storage: &pluginReportStorage{},
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
	/*	var options internalversion.ListOptions
		query := request.RequestQueryFrom(ctx)
		err := scheme.ParameterCodec.DecodeParameters(query, v1beta1.SchemeGroupVersion, &options)
		if nil != err {
			return err
		}*/

	klog.Infof("Get called with key: %v on resource %v\n", key, p.resourceName)

	split := strings.Split(key, "-")
	timestampStr := split[len(split)-1]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if nil != err {
		return fmt.Errorf("failed to parse timestampt %s, error: %w", timestampStr, err)
	}
	timeStr := time.Unix(timestamp, 0).Format(time.RFC3339)

	dir := "/report"
	readDir, err := os.ReadDir(dir)
	if nil != err {
		return err
	}
	var fileName string
	for _, item := range readDir {
		if item.IsDir() {
			continue
		}

		if strings.Contains(item.Name(), timeStr) {
			fileName = path.Join(dir, item.Name())
		}
	}
	if strings.EqualFold(fileName, "") {
		return fmt.Errorf("no task found")
	}

	file, err := os.Open(fileName)
	if nil != err {
		return err
	}
	readAll, err := io.ReadAll(file)
	if nil != err {
		return err
	}

	pluginReport := objPtr.(*v1beta1.PluginReport)
	err = json.Unmarshal(readAll, &(pluginReport.Spec))
	if nil != err {
		return err
	}
	pluginReport.Name = key
	pluginReport.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   v1beta1.GroupName,
		Version: "v1beta1",
		Kind:    "PluginReport",
	})

	return nil
}

func (p pluginReportStorage) GetList(ctx context.Context, key string, opts storage.ListOptions, listObj runtime.Object) error {
	dir := "/report"

	readDir, err := os.ReadDir(dir)
	if nil != err {
		return fmt.Errorf("failed to read directory %s, error: %w", dir, err)
	}

	pluginReportList := listObj.(*v1beta1.PluginReportList)
	var resList []runtime.Object
	for _, item := range readDir {
		if item.IsDir() {
			continue
		}

		fileName := path.Join(dir, item.Name())
		file, err := os.Open(fileName)
		if nil != err {
			return fmt.Errorf("failed to open file %s, error: %w", fileName, err)
		}
		readAll, err := io.ReadAll(file)
		if nil != err {
			return fmt.Errorf("failed to read file %s, error: %w", fileName, err)
		}

		pluginReport := &v1beta1.PluginReport{}
		err = json.Unmarshal(readAll, &(pluginReport.Spec))
		if nil != err {
			return fmt.Errorf("failed to unmarshal %#v into value %#v", readAll, pluginReport)
		}

		split := strings.Split(fileName, "_")
		timeStr := split[len(split)-1]
		times, err := time.Parse(time.RFC3339, timeStr)
		if nil != err {
			return fmt.Errorf("failed to parse time %s, error: %w", timeStr, err)
		}

		pluginReport.Name = fmt.Sprintf("%s-%d", pluginReport.Spec.TaskName, times.Unix())
		pluginReport.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
			Group:   v1beta1.GroupName,
			Version: "v1beta1",
			Kind:    "PluginReport",
		})

		resList = append(resList, pluginReport)
	}

	err = meta.SetList(pluginReportList, resList)
	if nil != err {
		return err
	}

	pluginReportList.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   v1beta1.GroupName,
		Version: "v1beta",
		Kind:    "PluginReportList",
	})

	return nil
}

func (p pluginReportStorage) GuaranteedUpdate(ctx context.Context, key string, destination runtime.Object, ignoreNotFound bool, preconditions *storage.Preconditions, tryUpdate storage.UpdateFunc, cachedExistingObject runtime.Object) error {
	return fmt.Errorf("GuaranteedUpdate API not implement")
}

func (p pluginReportStorage) Count(key string) (int64, error) {
	return 0, fmt.Errorf("Count not supported for key: %s", key)
}
