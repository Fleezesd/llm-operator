package v1alpha1

import (
	"context"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Finalizer is the name of the finalizer added to objects to ensure they are cleaned up.
	Finalizer     = Group + "/finalizer"
	ProviderLabel = Group + "/provider"
)

type ProviderType string

const (
	ProviderTypeUnknown  ProviderType = "unknown"
	ProviderType3rdParty ProviderType = "3rdParty"
	ProviderTypeWorker   ProviderType = "worker"
)

type Provider struct {
	Endpoint *Endpoint                    `json:"endpoint,omitempty"`
	Worker   *corev1.TypedObjectReference `json:"worker,omitempty"`
}

func (p Provider) GetType() ProviderType {
	if p.Endpoint != nil {
		return ProviderType3rdParty
	}
	if p.Worker != nil {
		return ProviderTypeWorker
	}
	return ProviderTypeUnknown
}

// Endpoint represents a reachable API endpoint.
type Endpoint struct {
	// URL for the endpoint.
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// InternalURL for this endpoint which is much faster but only can be used inside this cluster
	// +kubebuilder:validation:Required
	InternalURL string `json:"internalURL,omitempty"`

	// AuthSecret if the chart repository requires auth authentication,
	// set the username and password to secret, with the field user and password respectively.
	AuthSecret *corev1.TypedLocalObjectReference `json:"authSecret,omitempty"`

	// Insecure if the endpoint needs a secure connection
	Insecure bool `json:"insecure,omitempty"`
}

func (o Endpoint) AuthData(ctx context.Context, ns string, c client.Client) (map[string][]byte, error) {
	if lo.IsNil(o.AuthSecret) {
		return nil, nil
	}
	authSecret := &corev1.Secret{}
	if err := c.Get(ctx, types.NamespacedName{Name: o.AuthSecret.Name,
		Namespace: ns}, authSecret); err != nil {
		return nil, err
	}
	return authSecret.Data, nil
}

func (o Endpoint) AuthAPIKey(ctx context.Context, ns string, c client.Client) (string, error) {
	data, err := o.AuthData(ctx, ns, c)
	if err != nil {
		return "", err
	}
	return string(data["apiKey"]), nil
}
