package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	EventName        EventName
	PayloadInterface interface{}
}

func (e *TestEvent) Name() EventName {
	return e.EventName
}

func (e *TestEvent) DateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) Payload() interface{} {
	return e.Payload
}

type TestEventHandler struct{ ID int }

func (h *TestEventHandler) Handle(event EventInterface) {
}

// Suite test Event Dispatcher
type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = *NewEventDispatcher()
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.event = TestEvent{EventName: "test", PayloadInterface: "test"}
	suite.event2 = TestEvent{EventName: "test2", PayloadInterface: "test2"}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	// Register `handler` for `event`.
	err := suite.eventDispatcher.Register(suite.event.Name(), &suite.handler)
	// Test registration of `handler` went well.
	suite.Nil(err)
	// Test if the `handler` is in the list and unique, so far, for this `event`.
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.Name()]))

	// Register `handler2` for `event`.
	err = suite.eventDispatcher.Register(suite.event.Name(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.Name()]))

	// Test if the handlers were correctly registered (assigned).
	assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.handlers[suite.event.Name()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.Name()][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	// "Force" error situation; duplicate handler registration for the same event.
	// Test if the error is catched.
	err := suite.eventDispatcher.Register(suite.event.Name(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.Name()]))

	err = suite.eventDispatcher.Register(suite.event.Name(), &suite.handler)
	suite.Equal(ErrHandlerAlreadyRegistered, err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.Name()]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	err := suite.eventDispatcher.Register(suite.event.Name(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.Name()]))

	err = suite.eventDispatcher.Register(suite.event.Name(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.Name()]))

	err = suite.eventDispatcher.Register(suite.event2.Name(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.Name()]))

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	err := suite.eventDispatcher.Register(suite.event.Name(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.Name()]))

	err = suite.eventDispatcher.Register(suite.event.Name(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.Name()]))

	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.Name(), &suite.handler))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.Name(), &suite.handler2))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event.Name(), &suite.handler3))
	assert.False(suite.T(), suite.eventDispatcher.Has(EventName("non-existing-event"), &suite.handler))
}

type EventHandlerMock struct {
	mock.Mock
}

func (m *EventHandlerMock) Handle(event EventInterface) {
	m.Called(event)
}

func (suite *EventDispatcherTestSuite) TestEventDispatch_Dispatch() {
	eh := &EventHandlerMock{}
	eh.On("Handle", &suite.event)
	suite.eventDispatcher.Register(suite.event.Name(), eh)
	suite.eventDispatcher.Dispatch(&suite.event)
	// Make sure no errors happened, while making the Register or Dispatch calls.
	eh.AssertExpectations(suite.T())
	// Make sure `Handle` mock was called, once we dispatched an event.
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *EventDispatcherTestSuite) TestEventDispatch_Remove() {

	err := suite.eventDispatcher.Register(suite.event.Name(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.Name()]))

	err = suite.eventDispatcher.Register(suite.event.Name(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.Name()]))

	err = suite.eventDispatcher.Register(suite.event2.Name(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.Name()]))

	err = suite.eventDispatcher.Remove(suite.event.Name(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.Name()]))
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.Name()][0])
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.Name()]))
	assert.Equal(suite.T(), &suite.handler3, suite.eventDispatcher.handlers[suite.event2.Name()][0])

	err = suite.eventDispatcher.Remove(suite.event.Name(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event.Name()]))
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.Name()]))
	assert.Equal(suite.T(), &suite.handler3, suite.eventDispatcher.handlers[suite.event2.Name()][0])
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
