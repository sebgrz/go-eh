package goeh

import (
	"encoding/json"
	"fmt"
)

// Event abstraction
type Event interface {
	GetType() string
	GetPayload() string
	LoadPayload() error
	SavePayload(event Event) error
}

// EventData basic event structure which is use as base of all events
type EventData struct {
	Event
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// GetType of event
func (m *EventData) GetType() string {
	return m.Type
}

// GetPayload of event
func (m *EventData) GetPayload() string {
	return m.Payload
}

// LoadPayload : fetch event type from payload
func (m *EventData) LoadPayload() error {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(m.Payload), &data); err != nil {
		return err
	}
	m.Type = data["type"].(string)
	return nil
}

// SavePayload save event structure to payload json string
func (m *EventData) SavePayload(event Event) error {
	if m == nil {
		return fmt.Errorf("Event doesn't have base EventData instance. Please add this to struct instance -> EventData: new(archevents.EventData)")
	}

	m.Type = event.GetType()
	raw, err := json.Marshal(event)
	if err != nil {
		return err
	}

	m.Payload = string(raw)
	return nil
}
