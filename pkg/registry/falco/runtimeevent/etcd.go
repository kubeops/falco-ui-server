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

package request

import (
	api "kubeops.dev/falco-ui-server/apis/falco"
	"kubeops.dev/falco-ui-server/apis/falco/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

// ControllerStorage includes dummy storage for RuntimeEvents and for Status subresource.
type ControllerStorage struct {
	Controller *REST
}

func NewStorage(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*REST, error) {
	return NewREST(scheme, optsGetter)
}

type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against replication controllers.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*REST, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &api.RuntimeEvent{} },
		NewListFunc:              func() runtime.Object { return &api.RuntimeEventList{} },
		PredicateFunc:            MatchController,
		DefaultQualifiedResource: api.Resource(v1alpha1.ResourceRuntimeEvents),

		CreateStrategy:      strategy,
		UpdateStrategy:      strategy,
		DeleteStrategy:      strategy,
		ResetFieldsStrategy: strategy,

		TableConvertor: NewTableConvertor(api.Resource(v1alpha1.ResourceRuntimeEvents)),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	statusStrategy := statusStrategy{strategy: strategy}
	statusStore := *store
	statusStore.UpdateStrategy = statusStrategy
	statusStore.ResetFieldsStrategy = statusStrategy

	return &REST{store}, nil
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"rte"}
}

// Implement CategoriesProvider
var _ rest.CategoriesProvider = &REST{}

// Categories implements the CategoriesProvider interface. Returns a list of categories a resource is part of.
func (r *REST) Categories() []string {
	return []string{"falco"}
}
