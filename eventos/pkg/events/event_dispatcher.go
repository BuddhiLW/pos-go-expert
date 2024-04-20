package events

import "errors"

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventDispatcher struct {
	handlers map[EventName][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: map[EventName][]EventHandlerInterface{},
	}
}

func ExistHandler(ed *EventDispatcher, eventName EventName) bool {
	if _, exists := ed.handlers[eventName]; exists {
		return true
	} else {
		return false
	}
}

func SameHandler(ed *EventDispatcher, eventName EventName, handler EventHandlerInterface) error {
	for _, h := range ed.handlers[eventName] {
		if h == handler {
			return ErrHandlerAlreadyRegistered
		}
	}
	return nil
}

func (ed *EventDispatcher) Register(eventName EventName, handler EventHandlerInterface) error {
	if ExistHandler(ed, eventName) {
		err := SameHandler(ed, eventName, handler)
		if err != nil {
			return err
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}
