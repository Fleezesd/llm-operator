package v1alpha1

import (
	"context"

	"github.com/samber/lo"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (llm LLM) AuthAPIKey(ctx context.Context, c client.Client) (string, error) {
	if lo.IsNil(llm.Spec.Endpoint) {
		return "", nil
	}
	return llm.Spec.Endpoint.AuthAPIKey(ctx, llm.GetNamespace(), c)
}
