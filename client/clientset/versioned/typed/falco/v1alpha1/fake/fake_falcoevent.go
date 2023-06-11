/*
Copyright AppsCode Inc. and Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "kubeops.dev/falco-ui-server/apis/falco/v1alpha1"
)

// FakeFalcoEvents implements FalcoEventInterface
type FakeFalcoEvents struct {
	Fake *FakeFalcoV1alpha1
}

var falcoeventsResource = schema.GroupVersionResource{Group: "falco.appscode.com", Version: "v1alpha1", Resource: "falcoevents"}

var falcoeventsKind = schema.GroupVersionKind{Group: "falco.appscode.com", Version: "v1alpha1", Kind: "FalcoEvent"}

// Get takes name of the falcoEvent, and returns the corresponding falcoEvent object, and an error if there is any.
func (c *FakeFalcoEvents) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FalcoEvent, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(falcoeventsResource, name), &v1alpha1.FalcoEvent{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FalcoEvent), err
}

// List takes label and field selectors, and returns the list of FalcoEvents that match those selectors.
func (c *FakeFalcoEvents) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FalcoEventList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(falcoeventsResource, falcoeventsKind, opts), &v1alpha1.FalcoEventList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FalcoEventList{ListMeta: obj.(*v1alpha1.FalcoEventList).ListMeta}
	for _, item := range obj.(*v1alpha1.FalcoEventList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested falcoEvents.
func (c *FakeFalcoEvents) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(falcoeventsResource, opts))
}

// Create takes the representation of a falcoEvent and creates it.  Returns the server's representation of the falcoEvent, and an error, if there is any.
func (c *FakeFalcoEvents) Create(ctx context.Context, falcoEvent *v1alpha1.FalcoEvent, opts v1.CreateOptions) (result *v1alpha1.FalcoEvent, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(falcoeventsResource, falcoEvent), &v1alpha1.FalcoEvent{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FalcoEvent), err
}

// Update takes the representation of a falcoEvent and updates it. Returns the server's representation of the falcoEvent, and an error, if there is any.
func (c *FakeFalcoEvents) Update(ctx context.Context, falcoEvent *v1alpha1.FalcoEvent, opts v1.UpdateOptions) (result *v1alpha1.FalcoEvent, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(falcoeventsResource, falcoEvent), &v1alpha1.FalcoEvent{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FalcoEvent), err
}

// Delete takes name of the falcoEvent and deletes it. Returns an error if one occurs.
func (c *FakeFalcoEvents) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(falcoeventsResource, name, opts), &v1alpha1.FalcoEvent{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFalcoEvents) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(falcoeventsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FalcoEventList{})
	return err
}

// Patch applies the patch and returns the patched falcoEvent.
func (c *FakeFalcoEvents) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FalcoEvent, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(falcoeventsResource, name, pt, data, subresources...), &v1alpha1.FalcoEvent{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FalcoEvent), err
}