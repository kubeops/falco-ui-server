/*
Copyright AppsCode Inc. and Contributors

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

package openapi

import (
	"context"
	"strings"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
)

type RDStorage struct {
	cfg ResourceInfo
}

var (
	_ rest.GroupVersionKindProvider = &RDStorage{}
	_ rest.Scoper                   = &RDStorage{}
	_ rest.Lister                   = &RDStorage{}
	_ rest.Getter                   = &RDStorage{}
	_ rest.GracefulDeleter          = &RDStorage{}
	_ rest.Storage                  = &RDStorage{}
	_ rest.SingularNameProvider     = &RDStorage{}
)

func NewRDStorage(cfg ResourceInfo) *RDStorage {
	return &RDStorage{cfg}
}

func (r *RDStorage) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return r.cfg.gvk
}

func (r *RDStorage) GetSingularName() string {
	return strings.ToLower(r.cfg.gvk.Kind)
}

func (r *RDStorage) NamespaceScoped() bool {
	return r.cfg.namespaceScoped
}

// Getter
func (r *RDStorage) New() runtime.Object {
	return r.cfg.obj
}

func (r *RDStorage) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.New(), nil
}

// Lister
func (r *RDStorage) NewList() runtime.Object {
	return r.cfg.list
}

func (r *RDStorage) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	return r.NewList(), nil
}

func (r *RDStorage) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return new(metav1.Table), nil
}

// Deleter
func (r *RDStorage) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	return r.New(), true, nil
}

func (r *RDStorage) Destroy() {
}
