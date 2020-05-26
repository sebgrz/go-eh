# go-eh
An event handling library in Go

### Installation
`go get github.com/hetacode/go-eh`

### Example event implementation:
```golang
type UserEvent struct {
  *goeh.EventData
  FirstName string `json:"first_name"`
}

func (e *UserEvent) GetType() string {
  return "UserEvent"
}
```

This event should be register in the EventsMapper:

```golang
eventsMapper := new(goeh.EventsMapper)
eventsMapper.Register(new(UserEvent))
```

Now, we can resolve event instance by json string:

```golang
json := `{
  "type": "UserEvent",
  "first_name": "John"
}`

event, err := eventsMapper.Resolve(json)
if err != nil {
	t.Fatal(err)
}
  
trueEvent := event.(*UserEvent)
// trueEvent.FirstName == "John"
```

### EventHandler part

Event handler struct should implement an EventHandler interface with `Handle` function.

```golang
type UserEventHandler struct {
}

// Execute message
func (e *UserEventHandler) Handle(event goeh.Event) {
  e := event.(*UserEvent)
	/// some magic tricks here, processing and more...
}
```

EventHandler and his Event should be register in EventsHandlerManager:

```golang
manager := goeh.NewEventsHandlerManager()
if err := manager.Register(new(UserEvent), new(UserEventHandler); err != nil {
  t.Fatal(err)
}
```

In order to Execute appropriate Event, just call Execute function in EventsHandlerManager: 
```golang
if err := manager.Execute(event); err != nil {
	t.Fatal(err)
}
```

And the same with EventsMapper together:

```golang
json := `{
  "type": "UserEvent",
  "first_name": "John"
}`

event, err := eventsMapper.Resolve(json)
if err != nil {
	t.Fatal(err)
}

if err := manager.Execute(event); err != nil {
	t.Fatal(err)
}
```

That's all!
