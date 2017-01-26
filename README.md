go-events
===================


go-events is a small lightweight event bus for GoLang. It allows you to easily push events through your entire go application. 

----------

Installation
-------------
> go get github.com/jackhopner/go-events


Usage
-------------
Usage of this package is simple, simply instantiate a new event bus using:

`eventBus := bus.NewEventBus()`

Register instances to the bus using:

`eventBus.Register(instance)`

When you call Register all methods for that instance type (interface) which begin with 'Event_' will be registered, methods which begin with 'Event_Async_' will be registered and fired off in a go routine when being invoked.

Finally you can Submit events to your bus using:

`eventBus.Submit(&evt)`

All registered methods which take an argument that has the same type as the event passed will now be invoked with the event passed as an argument.

You can also call Deregister to remove an instance from the bus:

`eventBus.Deregister(instance)`

Here's a full example here:

```
package main

import (
	"log"

	bus "github.com/jackhopner/go-events"
)

type TestObject interface {
	Event_TestNormalEvent(evt *normalEvent)
}

type testObject struct{}

func (o *testObject) Event_TestNormalEvent(evt *normalEvent) {
	//do work here
	log.Printf("%+v", evt)
}

type normalEvent struct {
	message string
}

func newTestObject() TestObject {
	return &testObject{}
}

func main() {
	o := newTestObject()
	eventBus := bus.NewEventBus()
	eventBus.Register(o)

	evt := normalEvent{"test"}

	eventBus.Submit(&evt, false)
}
```

Output:
> 2017/01/26 21:41:12 &{message:test}

Notes
-------------
* Events should be submitted by reference
* Methods which receive events should handle all errors.
