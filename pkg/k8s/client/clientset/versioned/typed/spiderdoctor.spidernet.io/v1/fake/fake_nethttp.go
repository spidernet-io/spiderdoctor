// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	spiderdoctorspidernetiov1 "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeNethttps implements NethttpInterface
type FakeNethttps struct {
	Fake *FakeSpiderdoctorV1
}

var nethttpsResource = schema.GroupVersionResource{Group: "spiderdoctor.spidernet.io", Version: "v1", Resource: "nethttps"}

var nethttpsKind = schema.GroupVersionKind{Group: "spiderdoctor.spidernet.io", Version: "v1", Kind: "Nethttp"}

// Get takes name of the nethttp, and returns the corresponding nethttp object, and an error if there is any.
func (c *FakeNethttps) Get(ctx context.Context, name string, options v1.GetOptions) (result *spiderdoctorspidernetiov1.Nethttp, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(nethttpsResource, name), &spiderdoctorspidernetiov1.Nethttp{})
	if obj == nil {
		return nil, err
	}
	return obj.(*spiderdoctorspidernetiov1.Nethttp), err
}

// List takes label and field selectors, and returns the list of Nethttps that match those selectors.
func (c *FakeNethttps) List(ctx context.Context, opts v1.ListOptions) (result *spiderdoctorspidernetiov1.NethttpList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(nethttpsResource, nethttpsKind, opts), &spiderdoctorspidernetiov1.NethttpList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &spiderdoctorspidernetiov1.NethttpList{ListMeta: obj.(*spiderdoctorspidernetiov1.NethttpList).ListMeta}
	for _, item := range obj.(*spiderdoctorspidernetiov1.NethttpList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested nethttps.
func (c *FakeNethttps) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(nethttpsResource, opts))
}

// Create takes the representation of a nethttp and creates it.  Returns the server's representation of the nethttp, and an error, if there is any.
func (c *FakeNethttps) Create(ctx context.Context, nethttp *spiderdoctorspidernetiov1.Nethttp, opts v1.CreateOptions) (result *spiderdoctorspidernetiov1.Nethttp, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(nethttpsResource, nethttp), &spiderdoctorspidernetiov1.Nethttp{})
	if obj == nil {
		return nil, err
	}
	return obj.(*spiderdoctorspidernetiov1.Nethttp), err
}

// Update takes the representation of a nethttp and updates it. Returns the server's representation of the nethttp, and an error, if there is any.
func (c *FakeNethttps) Update(ctx context.Context, nethttp *spiderdoctorspidernetiov1.Nethttp, opts v1.UpdateOptions) (result *spiderdoctorspidernetiov1.Nethttp, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(nethttpsResource, nethttp), &spiderdoctorspidernetiov1.Nethttp{})
	if obj == nil {
		return nil, err
	}
	return obj.(*spiderdoctorspidernetiov1.Nethttp), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeNethttps) UpdateStatus(ctx context.Context, nethttp *spiderdoctorspidernetiov1.Nethttp, opts v1.UpdateOptions) (*spiderdoctorspidernetiov1.Nethttp, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(nethttpsResource, "status", nethttp), &spiderdoctorspidernetiov1.Nethttp{})
	if obj == nil {
		return nil, err
	}
	return obj.(*spiderdoctorspidernetiov1.Nethttp), err
}

// Delete takes name of the nethttp and deletes it. Returns an error if one occurs.
func (c *FakeNethttps) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(nethttpsResource, name, opts), &spiderdoctorspidernetiov1.Nethttp{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNethttps) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(nethttpsResource, listOpts)

	_, err := c.Fake.Invokes(action, &spiderdoctorspidernetiov1.NethttpList{})
	return err
}

// Patch applies the patch and returns the patched nethttp.
func (c *FakeNethttps) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *spiderdoctorspidernetiov1.Nethttp, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(nethttpsResource, name, pt, data, subresources...), &spiderdoctorspidernetiov1.Nethttp{})
	if obj == nil {
		return nil, err
	}
	return obj.(*spiderdoctorspidernetiov1.Nethttp), err
}
