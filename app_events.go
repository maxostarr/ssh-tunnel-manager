package main

import (
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

// Emit emits an event with the given name and data, returning a unique event ID.
func (a *App) Emit(eventName string, data ...interface{}) (string, error) {
	eventData := EventData{
		ID:   getId(),
		Data: data,
	}
	runtime.EventsEmit(a.ctx, eventName, eventData)
	return eventData.ID, nil
}

// EmitAndWait emits an event and waits for a response, returning the response data.
func (a *App) EmitAndWait(eventName string, data ...interface{}) (interface{}, error) {
	id, err := a.Emit(eventName, data)
	if err != nil {
		return nil, err
	}

	responseChannel := make(chan interface{})
	// It's assumed runtime.EventsOn handles registration in a way that doesn't block or can be done in a separate goroutine if necessary.
	// Error handling for EventsOn is omitted for brevity but should be considered.
	unregister := runtime.EventsOn(a.ctx, eventName, func(data ...interface{}) {
		// Type assertion with check
		eventData, ok := data[0].(EventData)
		if !ok {
			return
		}

		if eventData.ID == id {
			responseChannel <- eventData.Data
		}
	})

	// Unregister the event listener after the response is received
	defer unregister()

	return <-responseChannel, nil
}
