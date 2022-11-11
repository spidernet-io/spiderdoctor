// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	spiderdoctorspidernetiov1 "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	versioned "github.com/spidernet-io/spiderdoctor/pkg/k8s/client/clientset/versioned"
	internalinterfaces "github.com/spidernet-io/spiderdoctor/pkg/k8s/client/informers/externalversions/internalinterfaces"
	v1 "github.com/spidernet-io/spiderdoctor/pkg/k8s/client/listers/spiderdoctor.spidernet.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// NetdoctorInformer provides access to a shared informer and lister for
// Netdoctors.
type NetdoctorInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.NetdoctorLister
}

type netdoctorInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewNetdoctorInformer constructs a new informer for Netdoctor type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewNetdoctorInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredNetdoctorInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredNetdoctorInformer constructs a new informer for Netdoctor type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredNetdoctorInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SpiderdoctorV1().Netdoctors().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SpiderdoctorV1().Netdoctors().Watch(context.TODO(), options)
			},
		},
		&spiderdoctorspidernetiov1.Netdoctor{},
		resyncPeriod,
		indexers,
	)
}

func (f *netdoctorInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredNetdoctorInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *netdoctorInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&spiderdoctorspidernetiov1.Netdoctor{}, f.defaultInformer)
}

func (f *netdoctorInformer) Lister() v1.NetdoctorLister {
	return v1.NewNetdoctorLister(f.Informer().GetIndexer())
}
