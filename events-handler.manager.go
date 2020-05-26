package goeh

import (
	"fmt"
)

// EventsHandlerManager keep all bindings event <-> eventHandler
type EventsHandlerManager struct {
	handlers map[string]EventHandler
}

// NewEventsHandlerManager create new instance
func NewEventsHandlerManager() *EventsHandlerManager {
	manager := &EventsHandlerManager{
		handlers: make(map[string]EventHandler),
	}
	return manager
}

// Register commandHandler
// event - should by instacne of Event
// handler - just EventHandler
func (h *EventsHandlerManager) Register(event Event, handler EventHandler) error {
	eventTypeName := event.GetType()
	h.handlers[eventTypeName] = handler
	return nil
}

// Execute : find eventhandler and then execute an event
func (h *EventsHandlerManager) Execute(event Event) error {
	if h.handlers[event.GetType()] == nil {
		return fmt.Errorf("Cannot find event handler for %s event typ", event.GetType())
	}
	h.handlers[event.GetType()].Handle(event)

	return nil
}
