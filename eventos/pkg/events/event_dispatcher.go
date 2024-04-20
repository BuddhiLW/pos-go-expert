package events

import "errors"

var ErrHandlerAlreadyRegistered = errors.New("Handler already registered")
var ErrEventName = errors.New("Event not registered in the dispatcher")
var ErrDoesntHaveHandler = errors.New("Doesn't have this handler attached to the event")

type EventDispatcher struct {
	handlers map[EventName][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: map[EventName][]EventHandlerInterface{},
	}
}

func ExistEvent(ed *EventDispatcher, eventName EventName) bool {
	if _, exists := ed.handlers[eventName]; exists {
		return true
	} else {
		return false
	}
}

func hasHandler(ed *EventDispatcher, eventName EventName, handler EventHandlerInterface) bool {
	for _, h := range ed.handlers[eventName] {
		if h == handler {
			return true
		}
	}
	return false
}

func SameHandler(ed *EventDispatcher, eventName EventName, handler EventHandlerInterface) error {
	if hasHandler(ed, eventName, handler) {
		return ErrHandlerAlreadyRegistered
	}
	return nil
}

func (ed *EventDispatcher) Has(eventName EventName, handler EventHandlerInterface) bool {
	if ExistEvent(ed, eventName) && hasHandler(ed, eventName, handler) {
		return true
	}
	return false
}

func (ed *EventDispatcher) Register(eventName EventName, handler EventHandlerInterface) error {
	if ed.Has(eventName, handler) {
		return ErrHandlerAlreadyRegistered
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Clear() {
	ed.handlers = map[EventName][]EventHandlerInterface{}
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if !ExistEvent(ed, event.Name()) {
		return ErrEventName
	}

	// Dispatch all handlers `h`, associated with the event called `event.Name()`
	for _, h := range ed.handlers[event.Name()] {
		h.Handle(event)
	}
	return nil
}

func (ed *EventDispatcher) Remove(eventName EventName, handler EventHandlerInterface) error {
	if !ed.Has(eventName, handler) {
		return ErrDoesntHaveHandler
	}

	for i, h := range ed.handlers[eventName] {
		if h == handler {
			ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
			return nil
		}
	}
	return nil
}
