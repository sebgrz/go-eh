package goehtests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	goeh "github.com/hetacode/go-eh"
)

func TestCorrectRegisterEventHandlerWithEvent(t *testing.T) {
	manager := goeh.NewEventsHandlerManager()

	if err := manager.Register(new(TestEvent), new(TestEventHandler)); err != nil {
		t.Fatal(err)
	}
}

func TestSaveEventToPayload(t *testing.T) {
	event := &TestTwoEvent{
		EventData: new(goeh.EventData),
		FirstName: "Janusz",
	}

	if err := event.SavePayload(event); err != nil {
		t.Fatal(err)
	}
	payload := event.Payload
	assert.Equal(t, strings.Contains(payload, "TestTwoEvent"), true)
	assert.Equal(t, strings.Contains(payload, "first_name"), true)
	assert.Equal(t, strings.Contains(payload, "Janusz"), true)
}

func TestSaveEventToPayloadBadInitalizationEvent(t *testing.T) {
	event := &TestTwoEvent{
		// EventData: new(archevents.EventData),
		FirstName: "Janusz",
	}

	if err := event.SavePayload(event); err != nil {
		if strings.HasPrefix(err.Error(), "Event doesn't have base EventData instance") {
			return
		}
	}
	t.Fatal()
}

func TestCorrectRegisterAndExecuteEventHandlerWithEvent(t *testing.T) {
	manager := goeh.NewEventsHandlerManager()

	handler1 := new(TestEventHandler)
	if err := manager.Register(new(TestEvent), handler1); err != nil {
		t.Fatal(err)
	}
	handler2 := new(TestTwoEventHandler)
	if err := manager.Register(new(TestTwoEvent), handler2); err != nil {
		t.Fatal(err)
	}

	if err := manager.Execute(&TestTwoEvent{EventData: &goeh.EventData{ID: "abc123"}}); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, handler1.EventType, "")
	assert.Equal(t, handler2.EventType, "TestTwoEvent")
	assert.Equal(t, handler2.EventID, "abc123")
}

func TestResolveEventFromGivenJsonData(t *testing.T) {
	json := `{
		"id": "abc123",
		"type": "TestTwoEvent",
		"first_name": "Janusz"
	}`

	// Get event from json
	eventsMapper := new(goeh.EventsMapper)
	eventsMapper.Register(new(TestTwoEvent))
	eventsMapper.Register(new(TestEvent))

	event, err := eventsMapper.Resolve(json)
	if err != nil {
		t.Fatal(err)
	}

	// Check correct event type
	assert.Equal(t, event.GetID(), "abc123")
	assert.Equal(t, event.GetType(), "TestTwoEvent")
	assert.NotEqual(t, fmt.Sprintf("%T", event), fmt.Sprintf("%T", new(TestEvent)))
	assert.Equal(t, fmt.Sprintf("%T", event), fmt.Sprintf("%T", new(TestTwoEvent)))

	// Check FirstName value
	assert.Equal(t, event.(*TestTwoEvent).FirstName, "Janusz")
}

func TestResolveSameTwoEventsShouldHaveDifferentInstance(t *testing.T) {
	json1 := `{
		"id": "abc111",
		"type": "TestTwoEvent",
		"first_name": "Janusz"
	}`

	json2 := `{
		"id": "abc222",
		"type": "TestTwoEvent",
		"first_name": "Waldus"
	}`

	eventsMapper := new(goeh.EventsMapper)
	eventsMapper.Register(new(TestTwoEvent))

	event1, err := eventsMapper.Resolve(json1)
	if err != nil {
		t.Fatal(err)
	}
	event2, err := eventsMapper.Resolve(json2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, event1.GetType(), "TestTwoEvent")
	assert.Equal(t, event2.GetType(), "TestTwoEvent")
	assert.Equal(t, event1.(*TestTwoEvent).FirstName, "Janusz")
	assert.Equal(t, event2.(*TestTwoEvent).FirstName, "Waldus")
	if event1 == event2 {
		t.Fatalf("Events pointers are same: %p - %p", event1, event2)
	}
}

func TestResolveGivenJsonDataAsEventAndCallProperlyEventHandler(t *testing.T) {
	json := `{
		"id": "abc123",
		"type": "TestTwoEvent",
		"first_name": "Janusz"
	}`

	// Register Event Handlers
	manager := goeh.NewEventsHandlerManager()
	handler1 := new(TestEventHandler)
	if err := manager.Register(new(TestEvent), handler1); err != nil {
		t.Fatal(err)
	}
	handler2 := new(TestTwoEventHandler)
	if err := manager.Register(new(TestTwoEvent), handler2); err != nil {
		t.Fatal(err)
	}

	// Register Events
	eventsMapper := new(goeh.EventsMapper)
	eventsMapper.Register(new(TestTwoEvent))
	eventsMapper.Register(new(TestEvent))

	event, err := eventsMapper.Resolve(json)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, handler2.UserFirstName, "")
	// Execute properly event handler by passed event
	if err := manager.Execute(event); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, handler2.UserFirstName, "Janusz")
	assert.Equal(t, handler2.EventID, "abc123")
}

// TEST DATA
// #############################

// TestEvent to test eventHandler
type TestEvent struct {
	*goeh.EventData
}

func (e *TestEvent) GetType() string {
	return "TestEvent"
}

// TestTwoEvent : second event to test
type TestTwoEvent struct {
	*goeh.EventData
	FirstName string `json:"first_name"`
}

func (e *TestTwoEvent) GetType() string {
	return "TestTwoEvent"
}

// ##################

// TestEventHandler to test correct execution of messages
type TestEventHandler struct {
	EventType string
}

// Execute message
func (e *TestEventHandler) Handle(event goeh.Event) {
	e.EventType = event.GetType()
}

// TestTwoEventHandler to test correct execution of messages
type TestTwoEventHandler struct {
	EventID       string
	EventType     string
	UserFirstName string
}

// Execute message
func (e *TestTwoEventHandler) Handle(event goeh.Event) {
	e.EventID = event.GetID()
	e.EventType = event.GetType()
	e.UserFirstName = event.(*TestTwoEvent).FirstName
}
