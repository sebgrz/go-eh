package goeh

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// EventsMapper keep all registered events with their keys
type EventsMapper struct {
	events map[string]interface{}
}

// Register an event as pair string - event
func (m *EventsMapper) Register(event Event) {
	if m.events == nil {
		m.events = make(map[string]interface{})
	}
	eventType := event.GetType()
	m.events[eventType] = event
}

// Resolve and find event by json data
func (m *EventsMapper) Resolve(jsonData string) (Event, error) {
	event := &EventData{
		Payload: jsonData,
	}
	if err := event.LoadPayload(); err != nil {
		return nil, err
	}

	selectedEvent, ok := m.events[event.GetType()]
	if !ok {
		return nil, fmt.Errorf("Cannot find specific event: %s", event.GetType())
	}
	t := reflect.TypeOf(selectedEvent).Elem() // get value where pointer points to
	e := reflect.New(t).Interface()           // create a new instance
	if err := json.Unmarshal([]byte(jsonData), e); err != nil {
		return nil, err
	}

	return e.(Event), nil
}
