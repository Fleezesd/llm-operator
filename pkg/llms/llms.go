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

var (
	DefaultOpenAIModel string = "gpt-3.5-turbo"
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

type ModelParams interface {
	Marshal() []byte
	Unmarshal([]byte) error
}
