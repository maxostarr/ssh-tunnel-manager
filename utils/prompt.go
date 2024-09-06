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

// type PromptResponse []string

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
	fmt.Println("Prompt response data: ", data)
	if err != nil {
		return DefaultPromptResponse, err
	}

	response, ok := data.()
	if !ok {
		return DefaultPromptResponse, fmt.Errorf("invalid response data")
	}

	return response, nil
}
