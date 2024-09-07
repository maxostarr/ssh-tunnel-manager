package utils

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func getId() string {
	return uuid.New().String()
}

type EventData struct {
	ID   string
	Data interface{}
}

type EventManager interface {
	Emit(eventName string, data interface{}) (string, error)
	EmitAndWait(eventName string, data interface{}) ([]interface{}, error)
	Prompt(options PromptOptions) (PromptResponse, error)
}

type EventManagerImpl struct {
	ctx context.Context
}

func NewEventManager(ctx context.Context) EventManager {
	return EventManagerImpl{
		ctx: ctx,
	}
}

// Emit emits an event with the given name and data, returning a unique event ID.
func (m EventManagerImpl) Emit(eventName string, data interface{}) (string, error) {
	eventData := EventData{
		ID:   getId(),
		Data: data,
	}
	runtime.EventsEmit(m.ctx, eventName, eventData)
	return eventData.ID, nil
}

// EmitAndWait emits an event and waits for a response, returning the response data.
func (m EventManagerImpl) EmitAndWait(eventName string, data interface{}) ([]interface{}, error) {
	id, err := m.Emit(eventName, data)
	if err != nil {
		return nil, err
	}

	responseChannel := make(chan []interface{})
	responseEventName := eventName + id
	fmt.Println("Waiting for response event: " + responseEventName)
	// It's assumed runtime.EventsOn handles registration in a way that doesn't block or can be done in a separate goroutine if necessary.
	// Error handling for EventsOn is omitted for brevity but should be considered.
	runtime.EventsOnce(m.ctx, responseEventName, func(data ...interface{}) {
		// fmt.Println("Received response event: " + responseEventName)

		// // debug print data
		fmt.Println(data)

		// // Type assertion with check
		// eventData, ok := data[0].(EventData)
		// fmt.Println(eventData)

		// if !ok {
		// 	return
		// }

		// if eventData.ID == id {
		responseChannel <- data
		// }
	})

	// Unregister the event listener after the response is received
	// defer unregister()

	return <-responseChannel, nil
}
