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

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	llmv1alpha1 "github.com/fleezesd/llm-operator/api/v1alpha1"
	"github.com/fleezesd/llm-operator/pkg/model"
	"github.com/fleezesd/llm-operator/pkg/operator"
)

// ModelReconciler reconciles a Model object
type ModelReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=llm.fleezesd.io,resources=models,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=llm.fleezesd.io,resources=models/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=llm.fleezesd.io,resources=models/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=storageclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

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
	var m llmv1alpha1.Model

	err := r.Get(ctx, req.NamespacedName, &m)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	ctx = model.WithWrappedRecorder[*llmv1alpha1.Model](ctx, model.NewWrappedRecorder[*llmv1alpha1.Model](r.Recorder, &m))
	ctx = model.WithClient(ctx, r.Client)

	res, err := operator.ResultFromError(r.reconcile(ctx, req, &m))
	return operator.HandleError(ctx, res, err)
}

// reconcile model logic
func (r *ModelReconciler) reconcile(ctx context.Context, req ctrl.Request, m *llmv1alpha1.Model) error {
	client := model.ClientFromContext(ctx)
	recorder := model.WrappedRecorderFromContext[*llmv1alpha1.Model](ctx)

	// mean no available model need retry
	if !model.IsAvailable(ctx, *m) {
		hasSet, err := model.SetProgressing(ctx, client, *m)
		if err != nil {
			return err
		}
		if hasSet {
			recorder.Eventf("Normal", "ModelProgressing", "Model is progressing")
			return operator.RequeueAfter(time.Second)
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ModelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&llmv1alpha1.Model{}).
		Complete(r)
}

func (r *ModelReconciler) reconcilePVC(ctx context.Context, namespace string, name string, m *llmv1alpha1.Model) error {
	return nil
}
