package falcosidekick

import (
	"context"

	core "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct{}

func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&core.Pod{}).
		Complete(r)
}
