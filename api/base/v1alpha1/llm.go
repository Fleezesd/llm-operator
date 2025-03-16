package v1alpha1

import (
	"context"

	"github.com/fleezesd/llm-operator/pkg/llms"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (llm LLM) AuthAPIKey(ctx context.Context, c client.Client) (string, error) {
	if lo.IsNil(llm.Spec.Endpoint) {
		return "", nil
	}
	return llm.Spec.Endpoint.AuthAPIKey(ctx, llm.GetNamespace(), c)
}

func (llm LLM) Get3rdPartyModels() []string {
	if llm.Spec.Provider.GetType() != ProviderType3rdParty {
		return nil
	}

	if llm.Spec.Models != nil && len(llm.Spec.Models) > 0 {
		return llm.Spec.Models
	}

	switch llm.Spec.Type {
	case llms.OpenAI:
		return llms.OpenAIModels
	}
	return []string{}
}

// llm condition
func (llm LLM) ErrorCondition(msg string) Condition {
	currCon := llm.Status.GetCondition(TypeReady)

	if currCon.Type == TypeReady && currCon.Status == corev1.ConditionFalse && currCon.Message == msg {
		return currCon
	}

	lastSuccessfulTime := metav1.Now()
	if currCon.LastSuccessfulTime.IsZero() {
		lastSuccessfulTime = currCon.LastSuccessfulTime
	}
	return Condition{
		Type:               TypeReady,
		Status:             corev1.ConditionFalse,
		Reason:             ReasonAvailable,
		Message:            msg,
		LastSuccessfulTime: lastSuccessfulTime,
		LastTransitionTime: metav1.Now(),
	}
}

func (llm LLM) ReadyCondition(msg string) Condition {
	currCon := llm.Status.GetCondition(TypeReady)
	if currCon.Status == corev1.ConditionTrue && currCon.Reason == ReasonAvailable && currCon.Message == msg {
		return currCon
	}
	return Condition{
		Type:               TypeReady,
		Status:             corev1.ConditionTrue,
		Reason:             ReasonAvailable,
		Message:            msg,
		LastTransitionTime: metav1.Now(),
		LastSuccessfulTime: metav1.Now(),
	}
}
