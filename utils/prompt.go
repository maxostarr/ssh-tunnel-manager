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
	// Type is only allowed to be 'text' or 'password'
	Type PromptInputType
}

type PromptOptions struct {
	ConfirmText string
	CancelText  string
	Label       string
	Inputs      []PromptInput
}

type PromptResponse struct {
	Status   string
	Response []string
}

// Default empty prompt response
var DefaultPromptResponse = PromptResponse{
	Status:   "cancelled",
	Response: nil,
}

func (m EventManagerImpl) Prompt(options PromptOptions) (PromptResponse, error) {
	data, err := m.EmitAndWait("prompt", options)
	if err != nil {
		return DefaultPromptResponse, err
	}

	responseData, ok := data.([]interface{})
	if !ok {
		return DefaultPromptResponse, fmt.Errorf("invalid response data")
	}

	response := PromptResponse{
		Status:   responseData[0].(string),
		Response: responseData[1].([]string),
	}

	return response, nil
}
