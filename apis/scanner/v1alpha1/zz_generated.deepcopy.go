//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScanImage) DeepCopyInto(out *ScanImage) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Request != nil {
		in, out := &in.Request, &out.Request
		*out = new(ScanImageRequest)
		**out = **in
	}
	if in.Response != nil {
		in, out := &in.Response, &out.Response
		*out = new(ScanImageResponse)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScanImage.
func (in *ScanImage) DeepCopy() *ScanImage {
	if in == nil {
		return nil
	}
	out := new(ScanImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ScanImage) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScanImageRequest) DeepCopyInto(out *ScanImageRequest) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScanImageRequest.
func (in *ScanImageRequest) DeepCopy() *ScanImageRequest {
	if in == nil {
		return nil
	}
	out := new(ScanImageRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScanImageResponse) DeepCopyInto(out *ScanImageResponse) {
	*out = *in
	in.Result.DeepCopyInto(&out.Result)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScanImageResponse.
func (in *ScanImageResponse) DeepCopy() *ScanImageResponse {
	if in == nil {
		return nil
	}
	out := new(ScanImageResponse)
	in.DeepCopyInto(out)
	return out
}
