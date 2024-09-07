package utils

import (
	"fmt"
)

type PromptInputType string

const (
	PromptInputTypeText     PromptInputType = "text"
	PromptInputTypePassword PromptInputType = "password"
)

type PromptInput struct {
	Label string
	Key   string
	Type  PromptInputType
}

type PromptOptions struct {
	Type        string
	ConfirmText string
	CancelText  string
	Label       string
	Inputs      []PromptInput
}

type PromptResponse struct {
	Status   string
	Response map[string]string
}

// Default empty prompt response
var DefaultPromptResponse = PromptResponse{
	Status:   "cancelled",
	Response: nil,
}

func NewPromptOptions(label string, confirmText string, cancelText string, inputs []PromptInput) PromptOptions {
	return PromptOptions{
		Type:        "prompt",
		Label:       label,
		ConfirmText: confirmText,
		CancelText:  cancelText,
		Inputs:      inputs,
	}
}

func (m EventManagerImpl) Prompt(options PromptOptions) (PromptResponse, error) {
	data, err := m.EmitAndWait("prompt", options)
	if err != nil {
		return DefaultPromptResponse, err
	}

	status, ok := data[0].(string)
	if !ok {
		return DefaultPromptResponse, fmt.Errorf("invalid response status")
	}

	mapData, ok := data[1].(map[string]interface{})
	if !ok {
		return DefaultPromptResponse, fmt.Errorf("invalid response data")
	}

	response := make(map[string]string)

	for key, value := range mapData {
		response[key] = fmt.Sprintf("%v", value)
	}

	return PromptResponse{
		Status:   status,
		Response: response,
	}, nil

}
