package goeh

// EventHandler an abstraction to processing specific event
type EventHandler interface {
	Handle(event Event)
}
