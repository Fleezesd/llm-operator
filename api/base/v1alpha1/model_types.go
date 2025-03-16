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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ModelSpec defines the desired state of Model
type ModelSpec struct {
	CommonSpec `json:",inline"`

	// Type defines what kind of model this is
	// Comma separated field which can be wrapped by {llm,embedding}
	Types string `json:"types,omitempty"`

	// Source define the source of the model file
	Source *corev1.TypedObjectReference `json:"source,omitempty"`

	// HuggingFaceRepo defines the huggingface repo which hosts this model
	HuggingFaceRepo string `json:"huggingFaceRepo,omitempty"`
	// ModelScopeRepo defines th😍e modelscope repo which hosts this model
	ModelScopeRepo string `json:"modelScopeRepo,omitempty"`

	// Revision it's required if download model file from modelscope
	// It can be a tag, branch name.
	Revision string `json:"revision,omitempty"`

	ModelSource string `json:"modelSource,omitempty"`

	// MaxContextLength defines the max context length allowed in this model
	MaxContextLength int `json:"maxContextLength,omitempty"`
}

// ModelStatus defines the observed state of Model
type ModelStatus struct {
	// ConditionedStatus is the current status
	ConditionedStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Model is the Schema for the models API
type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelSpec   `json:"spec,omitempty"`
	Status ModelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ModelList contains a list of Model
type ModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Model `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Model{}, &ModelList{})
}
