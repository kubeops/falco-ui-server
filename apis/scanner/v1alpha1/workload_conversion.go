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

package v1alpha1

import (
	"fmt"

	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func (dst *Workload) Duckify(srcRaw runtime.Object) error {
	switch src := srcRaw.(type) {
	case *apps.Deployment:
		dst.TypeMeta = metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: apps.SchemeGroupVersion.String(),
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.Selector = src.Spec.Selector
		dst.Spec.Template = src.Spec.Template
		return nil
	case *apps.StatefulSet:
		dst.TypeMeta = metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: apps.SchemeGroupVersion.String(),
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.Selector = src.Spec.Selector
		dst.Spec.Template = src.Spec.Template
		return nil
	case *apps.DaemonSet:
		dst.TypeMeta = metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: apps.SchemeGroupVersion.String(),
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.Selector = src.Spec.Selector
		dst.Spec.Template = src.Spec.Template
		return nil
	case *unstructured.Unstructured:
		return runtime.DefaultUnstructuredConverter.FromUnstructured(src.UnstructuredContent(), dst)
	default:
		return fmt.Errorf("unknown src type %T", srcRaw)
	}
}