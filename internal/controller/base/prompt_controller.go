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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"errors"

	basev1alpha1 "github.com/fleezesd/llm-operator/api/base/v1alpha1"
	"github.com/fleezesd/llm-operator/pkg/llms"
	"github.com/fleezesd/llm-operator/pkg/llms/models/openai"
	"github.com/go-logr/logr"
	"github.com/samber/lo"
)

// PromptReconciler reconciles a Prompt object
type PromptReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=base.fleezesd.io,resources=prompts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=base.fleezesd.io,resources=prompts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=base.fleezesd.io,resources=prompts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Prompt object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PromptReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Starting prompt reconcile")
	// TODO(user): your logic here

	prompt := &basev1alpha1.Prompt{}
	if err := r.Get(ctx, req.NamespacedName, prompt); err != nil {
		logger.V(1).Info("unable to fetch prompt")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if newAdded := ctrlutil.AddFinalizer(prompt, basev1alpha1.Finalizer); newAdded {
		logger.Info("Try to add Finalizer for Prompt")
		if err := r.Update(ctx, prompt); err != nil {
			logger.Error(err, "Failed to update Prompt to add finalizer, will try again later")
			return ctrl.Result{}, err
		}
		logger.Info("Adding Finalizer for Prompt done")
		return ctrl.Result{Requeue: true}, nil
	}

	if prompt.GetDeletionTimestamp() != nil && ctrlutil.ContainsFinalizer(prompt, basev1alpha1.Finalizer) {
		logger.Info("Try to remove Finalizer for Prompt")
		ctrlutil.RemoveFinalizer(prompt, basev1alpha1.Finalizer)
		if err := r.Update(ctx, prompt); err != nil {
			logger.Error(err, "Failed to remove finalizer for Prompt")
			return ctrl.Result{}, err
		}
		logger.Info("Remove Prompt done")
		return ctrl.Result{}, nil
	}

	err := r.CallLLM(ctx, logger, prompt)
	if err != nil {
		logger.Error(err, "Failed to call llm")
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PromptReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&basev1alpha1.Prompt{}).
		Complete(r)
}

func (r *PromptReconciler) CallLLM(ctx context.Context, logger logr.Logger, prompt *basev1alpha1.Prompt) error {
	if err := r.validatePrompt(prompt); err != nil {
		return err
	}
	llm, err := r.getLLMConfig(ctx, prompt)
	if err != nil {
		return err
	}
	apiKey, err := llm.AuthAPIKey(ctx, r.Client)
	if err != nil {
		return r.UpdateStatus(ctx, prompt, nil, err)
	}
	llmClient, err := r.createLLMClient(llm, apiKey)
	if err != nil {
		return r.UpdateStatus(ctx, prompt, nil, err)
	}
	return r.executeLLMCall(ctx, prompt, llmClient)
}

func (r *PromptReconciler) validatePrompt(prompt *basev1alpha1.Prompt) error {
	if lo.IsNil(prompt.Spec.LLM) {
		return errors.New("no llm configuration provider")
	}
	return nil
}

func (r *PromptReconciler) getLLMConfig(ctx context.Context, prompt *basev1alpha1.Prompt) (*basev1alpha1.LLM, error) {
	llm := &basev1alpha1.LLM{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      prompt.Spec.LLM.Name,
		Namespace: prompt.Namespace,
	}, llm)
	return llm, err
}

func (r *PromptReconciler) createLLMClient(llm *basev1alpha1.LLM, apiKey string) (llms.LLM, error) {
	switch llm.Spec.Type {
	case llms.OpenAI:
		return openai.NewOpenAI(apiKey, llm.Spec.Endpoint.URL)
	default:
		return nil, errors.New("unsupported LLM type")
	}
}

func (r *PromptReconciler) executeLLMCall(ctx context.Context, prompt *basev1alpha1.Prompt, llmClient llms.LLM) error {
	var callData []byte
	// todo: need prompt marshal for call llm and add model params
	resp, err := llmClient.Call(callData)
	if err != nil {
		return r.UpdateStatus(ctx, prompt, resp, err)
	}
	return r.UpdateStatus(ctx, prompt, resp, nil)
}

func (r *PromptReconciler) UpdateStatus(ctx context.Context, prompt *basev1alpha1.Prompt,
	response llms.Response, err error) error {
	promptDeepCopy := prompt.DeepCopy()
	newCond := basev1alpha1.Condition{
		Type:               basev1alpha1.TypeDone,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             basev1alpha1.ReasonReconcileSuccess,
		Message:            "Finished",
	}
	if err != nil {
		newCond.Status = corev1.ConditionFalse
		newCond.Reason = basev1alpha1.ReasonReconcileError
		newCond.Message = err.Error()
	}
	promptDeepCopy.Status.SetConditions(newCond)
	if response != nil {
		promptDeepCopy.Status.Data = response.Bytes()
	}
	return errors.Join(err, r.Client.Status().Update(ctx, promptDeepCopy))
}
