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

	llmopenai "github.com/fleezesd/llm-operator/pkg/llms/models/openai"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PromptSpec defines the desired state of Prompt
type PromptSpec struct {
	// llm service name (CRD LLM)
	LLM *corev1.TypedObjectReference `json:"llm"`
	// OpenAI Prompt Params
	OpenAIParams *llmopenai.ModelParams `json:"openAIParams,omitempty"`
}

// PromptStatus defines the observed state of Prompt
type PromptStatus struct {
	ConditionedStatus `json:",inline"`
	// Data retrieved after LLM Call
	Data []byte `json:"data"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Prompt is the Schema for the prompts API
type Prompt struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PromptSpec   `json:"spec,omitempty"`
	Status PromptStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PromptList contains a list of Prompt
type PromptList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Prompt `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Prompt{}, &PromptList{})
}
