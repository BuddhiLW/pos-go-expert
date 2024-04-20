package events

import "time"

type EventInterface interface {
	Name() string
	DateTime() time.Time
	Payload() interface{}
}

type EventHandlerInterface interface {
	Handle(event EventInterface)
}

type EventDispatcherInterface interface {
	// Dispatch(event EventInterface)
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear() error
}
