package events

import "time"

type EventName string

type EventInterface interface {
	Name() EventName
	DateTime() time.Time
	Payload() interface{}
}

type EventHandlerInterface interface {
	Handle(event EventInterface)
}

// Dispath(Event) -> []EventHandler <=> Event ~ []EventHandler

type EventDispatcherInterface interface {
	// Dispatch(event EventInterface)
	Register(eventName EventName, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(eventName EventName, handler EventHandlerInterface) error
	Has(eventName EventName, handler EventHandlerInterface) bool
	Clear() error
}
