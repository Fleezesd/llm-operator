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
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	basev1alpha1 "github.com/fleezesd/llm-operator/api/base/v1alpha1"
	"github.com/go-logr/logr"
)

// LLMReconciler reconciles a LLM object
type LLMReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=base.fleezesd.io,resources=llms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=base.fleezesd.io,resources=llms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=base.fleezesd.io,resources=llms/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LLM object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *LLMReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	logger := log.FromContext(ctx)
	logger.Info("Reconciling LLM resource")

	instance := &basev1alpha1.LLM{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		logger.V(1).Info("Failed to get LLM")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// add finalizer
	if newAdded := ctrlutil.AddFinalizer(instance, basev1alpha1.Finalizer); newAdded {
		logger.Info("Try to add finalizer for LLM")
		if err := r.Update(ctx, instance); err != nil {
			logger.Error(err, "Failed to update LLM to add finalizer, will try again later")
			return ctrl.Result{}, err
		}
		logger.Info("Adding Finalizer for LLM done")
		return ctrl.Result{Requeue: true}, nil
	}

	// check if the lllm instance is marked to be deleted
	if instance.GetDeletionTimestamp() != nil && ctrlutil.ContainsFinalizer(instance, basev1alpha1.Finalizer) {
		logger.Info("Performing Finalizer Operations for LLM before delete CR")
		// TODO perform the finalizer operations here, for example: remove data?
		logger.Info("Removing Finalizer for LLM after successfully performing the operations")
		ctrlutil.RemoveFinalizer(instance, basev1alpha1.Finalizer)
		if err := r.Update(ctx, instance); err != nil {
			logger.Error(err, "Failed to remove finalizer for LLM")
			return ctrl.Result{}, err
		}
		logger.Info("Remove LLM done")
		return ctrl.Result{}, nil
	}

	// check label
	if instance.Labels == nil {
		instance.Labels = make(map[string]string)
	}
	providerType := instance.Spec.Provider.GetType()
	if _type, ok := instance.Labels[basev1alpha1.ProviderLabel]; !ok || _type != string(providerType) {
		instance.Labels[basev1alpha1.ProviderLabel] = string(providerType)
		err := r.Client.Update(ctx, instance)
		if err != nil {
			logger.Error(err, "failed to update llm labels", "providerType", providerType)
		}
		return ctrl.Result{Requeue: true}, err
	}

	// check llm
	err := r.CheckLLM(ctx, logger, instance)
	if err != nil {
		logger.Error(err, "Failed to check LLM")
		// Update conditioned status
		return ctrl.Result{RequeueAfter: waitMedium}, err
	}

	return ctrl.Result{RequeueAfter: waitLonger}, nil
}

func (r *LLMReconciler) CheckLLM(ctx context.Context, logger logr.Logger, instance *basev1alpha1.LLM) error {
	logger.Info("Checking LLM resource")

	switch instance.Spec.Provider.GetType() {
	case basev1alpha1.ProviderType3rdParty:
		return r.check3rdPartyLLM(ctx, logger, instance)
	case basev1alpha1.ProviderTypeWorker:
		return r.checkWorkerLLM(ctx, logger, instance)
	}
	return nil
}

func (r *LLMReconciler) check3rdPartyLLM(ctx context.Context, logger logr.Logger, instance *basev1alpha1.LLM) error {
	logger.Info("Checking 3rd party LLM resource")

	var err error
	// var msg string

	_, err = instance.AuthAPIKey(ctx, r.Client)
	if err != nil {
		return r.UpdateStatus(ctx, instance, "", err)
	}
	// todo: check 3rd party llm

	return nil
}

func (r *LLMReconciler) UpdateStatus(ctx context.Context, instance *basev1alpha1.LLM, t interface{}, err error) error {
	instanceCopy := instance.DeepCopy()
	var newCondition basev1alpha1.Condition
	if err != nil {
		newCondition = instance.ErrorCondition(err.Error())
	} else {
		msg, ok := t.(string)
		if !ok {
			msg = statusNilResponse
		}
		newCondition = instance.ReadyCondition(msg)
	}
	instanceCopy.Status.SetConditions(newCondition)
	return errors.Join(err, r.Client.Status().Update(ctx, instanceCopy))
}

func (r *LLMReconciler) checkWorkerLLM(ctx context.Context, logger logr.Logger, instance *basev1alpha1.LLM) error {
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LLMReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&basev1alpha1.LLM{}).
		Complete(r)
}
