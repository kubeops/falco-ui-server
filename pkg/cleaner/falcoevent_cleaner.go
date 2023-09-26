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

package cleaner

import (
	"context"
	"time"

	api "kubeops.dev/falco-ui-server/apis/falco/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func StartCleaner(kc client.Client, ttl time.Duration) {
	klog.Infoln("Starts the FalcoEvent cleaner")
	for range time.Tick(30 * time.Minute) {
		err := cleanerFunc(kc, ttl)
		if err != nil {
			klog.Errorf("Error occurred while cleaning Falco Events : %s \n", err.Error())
		}
	}
}

func cleanerFunc(kc client.Client, ttl time.Duration) error {
	var evList api.FalcoEventList
	err := kc.List(context.TODO(), &evList)
	if err != nil {
		return err
	}
	for _, ev := range evList.Items {
		if time.Since(ev.Spec.Time.Time) >= ttl {
			err = kc.Delete(context.TODO(), &ev)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}
	}
	return nil
}
