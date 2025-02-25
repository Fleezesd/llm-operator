package model

import (
	"context"

	"github.com/fleezesd/llm-operator/api/v1alpha1"
	llmv1alpha1 "github.com/fleezesd/llm-operator/api/v1alpha1"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func IsAvailable(ctx context.Context, m llmv1alpha1.Model) bool {
	return len(lo.Filter(m.Status.Conditions,
		func(item v1alpha1.ModelStatusCondition, _ int) bool {
			return item.Type == v1alpha1.ModelAvailable
		},
	)) > 0
}

func SetProgressing(
	ctx context.Context,
	c client.Client,
	m llmv1alpha1.Model,
) (bool, error) {
	hasProgressing := len(lo.Filter(m.Status.Conditions, func(item llmv1alpha1.ModelStatusCondition, _ int) bool {
		return item.Type == llmv1alpha1.ModelProgressing
	})) > 0
	if hasProgressing {
		return false, nil
	}
	m.Status.Conditions = []llmv1alpha1.ModelStatusCondition{
		{
			Type:               llmv1alpha1.ModelProgressing,
			Status:             corev1.ConditionTrue,
			LastUpdateTime:     metav1.Now(),
			LastTransitionTime: metav1.Now(),
		},
	}
	err := c.Status().Update(ctx, &m)
	if err != nil {
		return false, err
	}
	return true, nil
}
