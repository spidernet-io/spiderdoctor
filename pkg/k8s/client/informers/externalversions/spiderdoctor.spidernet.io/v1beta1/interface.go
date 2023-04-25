// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	internalinterfaces "github.com/spidernet-io/spiderdoctor/pkg/k8s/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Netdnses returns a NetdnsInformer.
	Netdnses() NetdnsInformer
	// Nethttps returns a NethttpInformer.
	Nethttps() NethttpInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Netdnses returns a NetdnsInformer.
func (v *version) Netdnses() NetdnsInformer {
	return &netdnsInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Nethttps returns a NethttpInformer.
func (v *version) Nethttps() NethttpInformer {
	return &nethttpInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}