package model

import (
	"context"

	"github.com/fleezesd/llm-operator/api/v1alpha1"
	llmv1alpha1 "github.com/fleezesd/llm-operator/api/v1alpha1"
	"github.com/samber/lo"
)

func IsAvailable(ctx context.Context, m llmv1alpha1.Model) bool {
	availableModels := lo.Filter(m.Status.Conditions,
		func(item v1alpha1.ModelStatusCondition, _ int) bool {
			return item.Type == v1alpha1.ModelAvailable
		},
	)
	return len(availableModels) > 0
}
