package openai

import (
	"encoding/json"

	"github.com/fleezesd/llm-operator/pkg/llms"
)

type Response struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

func (response *Response) Type() llms.LLMType {
	return llms.OpenAI
}

func (response *Response) Bytes() []byte {
	bytes, err := json.Marshal(response)
	if err != nil {
		return []byte{}
	}
	return bytes
}

func (response *Response) String() string {
	return string(response.Bytes())
}

func (response *Response) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, response)
}
