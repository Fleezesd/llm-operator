package openai

import (
	"context"
	"time"

	"github.com/fleezesd/llm-operator/pkg/llms"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	langchainllms "github.com/tmc/langchaingo/llms"
	langchainopenai "github.com/tmc/langchaingo/llms/openai"
)

const (
	OpenAIModelAPIURL    = "https://api.openai.com/v1"
	OpenAIDefaultTimeout = 300 * time.Second
)

var _ llms.LLM = (*OpenAI)(nil)

type OpenAI struct {
	apiKey  string
	baseURL string
}

func NewOpenAI(apiKey, baseURL string) (*OpenAI, error) {
	if apiKey == "" {
		return nil, errors.New("API key cannot be empty")
	}

	client := &OpenAI{
		apiKey:  apiKey,
		baseURL: lo.Ternary(baseURL == "", OpenAIModelAPIURL, baseURL),
	}

	return client, nil
}

func (o *OpenAI) Type() llms.LLMType {
	return llms.OpenAI
}

func (o *OpenAI) Call(input []byte) (llms.Response, error) {
	ctx := context.Background()
	// default use gpt-3.5 turbo
	llm, err := langchainopenai.New(
		langchainopenai.WithBaseURL(o.baseURL),
		langchainopenai.WithToken(o.apiKey),
	)
	if err != nil {
		return nil, errors.Errorf("init openai client: %v", err)
	}

	resp, err := llm.Call(ctx, string(input))
	if err != nil {
		return nil, err
	}

	return &Response{
		Code:    200,
		Data:    resp,
		Msg:     "call openai model success",
		Success: true,
	}, nil
}

func (o *OpenAI) Validate(ctx context.Context, options ...langchainllms.CallOption) (llms.Response, error) {
	llm, err := langchainopenai.New(
		langchainopenai.WithBaseURL(o.baseURL),
		langchainopenai.WithToken(o.apiKey),
	)
	if err != nil {
		return nil, errors.Errorf("init openai client: %d", err)
	}

	resp, err := llm.Call(ctx, "Hello", options...)
	if err != nil {
		return nil, err
	}

	return &Response{
		Code:    200,
		Data:    resp,
		Msg:     "",
		Success: true,
	}, nil
}
