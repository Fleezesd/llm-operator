package openai

import (
	"encoding/json"

	"github.com/fleezesd/llm-operator/pkg/llms"
)

type Role string

const (
	User      Role = "user"
	Assistant Role = "assistant"
)

var _ llms.ModelParams = (*ModelParams)(nil)

// +kubebuilder:object:generate=true
type ModelParams struct {
	// Method used for this prompt call
	Method Method `json:"method,omitempty"`

	// Model used for this prompt call
	Model string `json:"model,omitempty"`

	// Temperature is float in openai
	Temperature float32 `json:"temperature,omitempty"`

	// TopP is float in openai
	TopP float32 `json:"top_p,omitempty"`

	// Contents
	Prompt []Prompt `json:"prompt"`

	// TaskID is used for getting result of AsyncInvoke
	TaskID string `json:"task_id,omitempty"`

	// Incremental is only Used for SSE Invoke
	Incremental bool `json:"incremental,omitempty"`
}

type Prompt struct {
	Role    Role   `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

func DefaultModelParams() ModelParams {
	return ModelParams{
		Model:       llms.DefaultOpenAIModel,
		Method:      OpenAIInvoke,
		Temperature: 0.8,
		TopP:        0.7,
		Prompt:      []Prompt{},
	}
}

func (params *ModelParams) Marshal() []byte {
	data, err := json.Marshal(params)
	if err != nil {
		return []byte{}
	}
	return data
}

func (params *ModelParams) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, params)
}
