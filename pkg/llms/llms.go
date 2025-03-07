package llms

import (
	"context"

	langchainllms "github.com/tmc/langchaingo/llms"
)

type LLMType string

const (
	OpenAI   LLMType = "openai"
	Deepseek LLMType = "deepseek"
)

type LLM interface {
	Type() LLMType
	Call([]byte) (Response, error)
	Validate(context.Context, ...langchainllms.CallOption) (Response, error)
}

type Response interface {
	Type() LLMType
	String() string
	Bytes() []byte
	Unmarshal([]byte) error
}
