package events

import (
	"testing"
	"time"
)

type eventBusFixture struct {
	eventBus   EventBus
	testObject TestObject
}

type TestObject interface {
	Event_TestNormalEvent(evt *normalEvent)
	Event_Async_TestAsyncEvent(evt *asyncEvent)
	GetLastNormalEvent() *normalEvent
	GetAsyncChan() chan *asyncEvent
}

type testObject struct {
	receivedNormal *normalEvent
	receivedAsync  chan *asyncEvent
}

func (o *testObject) GetLastNormalEvent() *normalEvent {
	return o.receivedNormal
}

func (o *testObject) GetAsyncChan() chan *asyncEvent {
	return o.receivedAsync
}

func (o *testObject) Event_TestNormalEvent(evt *normalEvent) {
	o.receivedNormal = evt
}

func (o *testObject) Event_Async_TestAsyncEvent(evt *asyncEvent) {
	o.receivedAsync <- evt
}

type normalEvent struct {
	message string
}

type asyncEvent struct {
	number int64
}

func Test_NormalEvent(t *testing.T) {
	f := initEventBusFixture()
	defer close(f.testObject.GetAsyncChan())

	evt := normalEvent{"test"}

	f.eventBus.Submit(&evt, false)
	lastEvt := f.testObject.GetLastNormalEvent()
	if lastEvt == nil || lastEvt.message != evt.message {
		t.Fatalf(
			"Expected last event to have message [%s] was [%s]",
			evt.message,
			lastEvt.message,
		)
	}
	asyncChan := f.testObject.GetAsyncChan()
	select {
	case m := <-asyncChan:
		t.Fatalf(
			"Expected no async event, got at least [%+v]",
			m,
		)
	case <-time.After(1 * time.Second):
	}
}

func Test_AsyncEvent(t *testing.T) {
	f := initEventBusFixture()
	defer close(f.testObject.GetAsyncChan())

	evt := asyncEvent{int64(33)}

	f.eventBus.Submit(&evt, false)
	lastEvt := f.testObject.GetLastNormalEvent()
	if lastEvt != nil {
		t.Fatalf(
			"Expected last event to be nil, was [%s]",
			evt,
		)
	}

	asyncChan := f.testObject.GetAsyncChan()
	select {
	case m := <-asyncChan:
		if m.number != evt.number {
			t.Fatalf(
				"Expected last event to have number [%d] was [%d]",
				evt.number,
				m.number,
			)
		}
	case <-time.After(1 * time.Second):
		t.Fatal(
			"Expected async event, got none :(",
		)
	}
}

func Test_Submit_AsyncEvent(t *testing.T) {
	f := initEventBusFixture()
	defer close(f.testObject.GetAsyncChan())

	evt := normalEvent{"test"}

	f.eventBus.Submit(&evt, true)
	lastEvt := f.testObject.GetLastNormalEvent()
	if lastEvt != nil {
		t.Fatalf(
			"Expected last event to be nil, was [%s]",
			evt,
		)
	}
	time.Sleep(1 * time.Second)
	lastEvt = f.testObject.GetLastNormalEvent()
	if lastEvt == nil || lastEvt.message != evt.message {
		t.Fatalf(
			"Expected last event to have message [%s] was [%s]",
			evt.message,
			lastEvt.message,
		)
	}

	asyncChan := f.testObject.GetAsyncChan()
	select {
	case m := <-asyncChan:
		t.Fatalf(
			"Expected no async event, got at least [%+v]",
			m,
		)
	case <-time.After(1 * time.Second):
	}
}

func newT() TestObject {
	return &testObject{nil, make(chan *asyncEvent)}
}

func initEventBusFixture() eventBusFixture {
	eventBus := NewEventBus()
	testObject := newT()
	eventBus.Register(testObject)

	return eventBusFixture{eventBus, testObject}
}
