package events

import (
	"reflect"
	"strings"
)

type EventBus interface {
	Submit(event interface{}, async bool)
	Register(instance interface{})
	Deregister(instance interface{})
}

type eventBus struct {
	typeToEventHandlers map[reflect.Type][]EventHandler
}

type EventHandler struct {
	instance interface{}
	method   reflect.Value
	types    []reflect.Type
	async    bool
}

func remove(slice []EventHandler, s int) []EventHandler {
	return append(slice[:s], slice[s+1:]...)
}

func createHandlers(instance interface{}) []EventHandler {
	instanceType := reflect.TypeOf(instance)
	methods := getMethods(instance)
	handlers := []EventHandler{}

	for idx, method := range methods {
		name := instanceType.Method(idx).Name
		if strings.HasPrefix(name, "Event_") {
			async := strings.HasPrefix(name, "Event_Async_")

			types := getTypes(method)
			if len(types) == 0 {
				continue
			}

			handler := EventHandler{&instance, method, types, async}
			handlers = append(handlers, handler)
		}
	}

	return handlers
}

func invokeHandler(event interface{}, eventType reflect.Type, handler EventHandler) {
	args := make([]reflect.Value, len(handler.types))
	for i, typ := range handler.types {
		if typ == eventType {
			args[i] = reflect.ValueOf(event)
		} else {
			args[i] = reflect.Zero(typ)
		}
	}
	handler.method.Call(args)
}

func (bus *eventBus) Submit(event interface{}, async bool) {
	eventValue := reflect.ValueOf(event)
	eventType := eventValue.Type()
	if handlers, ok := bus.typeToEventHandlers[eventType]; ok {
		for _, h := range handlers {
			if async || h.async {
				go func(evt interface{}, handler EventHandler) {
					invokeHandler(evt, eventType, handler)
				}(event, h)
			} else {
				invokeHandler(event, eventType, h)
			}
		}
	}
}

func (bus *eventBus) Register(instance interface{}) {
	handlers := createHandlers(instance)
	for _, handler := range handlers {
		for _, typ := range handler.types {
			typeHandlers, ok := bus.typeToEventHandlers[typ]
			if !ok {
				typeHandlers = []EventHandler{}
				bus.typeToEventHandlers[typ] = typeHandlers
			}
			bus.typeToEventHandlers[typ] = append(typeHandlers, handler)
		}
	}
}

func (bus *eventBus) Deregister(instance interface{}) {
	instanceValue := reflect.ValueOf(instance)
	numMethods := instanceValue.NumMethod()

	for m := 0; m < numMethods; m++ {
		method := instanceValue.Method(m)

		types := getTypes(method)
		if len(types) == 0 {
			continue
		}

		for _, typ := range types {
			if typeHandlers, ok := bus.typeToEventHandlers[typ]; ok {
				toRemove := []int{}
				for i, handler := range typeHandlers {
					if handler.instance == instance {
						toRemove = append(toRemove, i)
					}
				}
				for _, removeIdx := range toRemove {
					typeHandlers = remove(typeHandlers, removeIdx)
				}
			}
		}
	}
}

func NewEventBus() EventBus {
	return &eventBus{map[reflect.Type][]EventHandler{}}
}
