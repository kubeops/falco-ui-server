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

package falcoevent

import (
	"context"
	"time"

	api "kubeops.dev/falco-ui-server/apis/falco/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// FalcoEventReconciler reconciles a FalcoEvent object
type FalcoEventReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	ReportTTL time.Duration
}

func (r *FalcoEventReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var isrp api.FalcoEvent
	if err := r.Get(ctx, req.NamespacedName, &isrp); err != nil {
		log.Error(err, "unable to fetch FalcoEvent")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if time.Since(isrp.Spec.Time.Time) >= r.ReportTTL {
		return ctrl.Result{}, r.Delete(ctx, &isrp)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FalcoEventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.FalcoEvent{}).
		Complete(r)
}
