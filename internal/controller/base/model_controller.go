/*
Copyright 2025.

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

package base

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	basev1alpha1 "github.com/fleezesd/llm-operator/api/base/v1alpha1"
	"github.com/go-logr/logr"
)

// ModelReconciler reconciles a Model object
type ModelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=base.fleezesd.io,resources=models,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=base.fleezesd.io,resources=models/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=base.fleezesd.io,resources=models/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Model object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *ModelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.V(5).Info("Starting model reconcile")

	model := &basev1alpha1.Model{}
	if err := r.Get(ctx, req.NamespacedName, model); err != nil {
		logger.V(1).Info("Failed to get Model")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if newAdded := ctrlutil.AddFinalizer(model, basev1alpha1.Finalizer); newAdded {
		logger.Info("Try to add Finalizer for Model")
		if err := r.Update(ctx, model); err != nil {
			logger.Error(err, "Failed to update Model to add finalizer, will try again later")
			return ctrl.Result{}, err
		}
		logger.Info("Adding Finalizer for Model done")
		return ctrl.Result{Requeue: true}, nil
	}

	if model.GetDeletionTimestamp() != nil && ctrlutil.ContainsFinalizer(model, basev1alpha1.Finalizer) {
		logger.Info("Performing Finalizer Operations for Model before delete CR")
		// remova all model files from storage service

	}

	return ctrl.Result{}, nil
}
func (r *ModelReconciler) RemoveModel(ctx context.Context, logger logr.Logger, model *basev1alpha1.Model) error {

	return nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *ModelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&basev1alpha1.Model{}).
		Complete(r)
}
