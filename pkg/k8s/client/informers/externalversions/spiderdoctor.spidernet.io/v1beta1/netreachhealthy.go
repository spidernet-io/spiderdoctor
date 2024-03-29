// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	spiderdoctorspidernetiov1beta1 "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
	versioned "github.com/spidernet-io/spiderdoctor/pkg/k8s/client/clientset/versioned"
	internalinterfaces "github.com/spidernet-io/spiderdoctor/pkg/k8s/client/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/spidernet-io/spiderdoctor/pkg/k8s/client/listers/spiderdoctor.spidernet.io/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// NetReachHealthyInformer provides access to a shared informer and lister for
// NetReachHealthies.
type NetReachHealthyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.NetReachHealthyLister
}

type netReachHealthyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewNetReachHealthyInformer constructs a new informer for NetReachHealthy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewNetReachHealthyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredNetReachHealthyInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredNetReachHealthyInformer constructs a new informer for NetReachHealthy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredNetReachHealthyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SpiderdoctorV1beta1().NetReachHealthies().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SpiderdoctorV1beta1().NetReachHealthies().Watch(context.TODO(), options)
			},
		},
		&spiderdoctorspidernetiov1beta1.NetReachHealthy{},
		resyncPeriod,
		indexers,
	)
}

func (f *netReachHealthyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredNetReachHealthyInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *netReachHealthyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&spiderdoctorspidernetiov1beta1.NetReachHealthy{}, f.defaultInformer)
}

func (f *netReachHealthyInformer) Lister() v1beta1.NetReachHealthyLister {
	return v1beta1.NewNetReachHealthyLister(f.Informer().GetIndexer())
}
