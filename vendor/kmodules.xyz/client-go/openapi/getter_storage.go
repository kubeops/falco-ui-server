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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
)

type GetterStorage struct {
	cfg ResourceInfo
}

var (
	_ rest.GroupVersionKindProvider = &GetterStorage{}
	_ rest.Scoper                   = &GetterStorage{}
	_ rest.Getter                   = &GetterStorage{}
	_ rest.Storage                  = &GetterStorage{}
	_ rest.SingularNameProvider     = &GetterStorage{}
)

func NewGetterStorage(cfg ResourceInfo) *GetterStorage {
	return &GetterStorage{cfg}
}

func (r *GetterStorage) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return r.cfg.gvk
}

func (r *GetterStorage) GetSingularName() string {
	return strings.ToLower(r.cfg.gvk.Kind)
}

func (r *GetterStorage) NamespaceScoped() bool {
	return r.cfg.namespaceScoped
}

// Getter
func (r *GetterStorage) New() runtime.Object {
	return r.cfg.obj
}

func (r *GetterStorage) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.New(), nil
}

func (r *GetterStorage) Destroy() {
}
