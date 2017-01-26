package events

import "reflect"

func getMethods(instance interface{}) map[int]reflect.Value {
	instanceValue := reflect.ValueOf(instance)
	numMethods := instanceValue.NumMethod()
	methods := map[int]reflect.Value{}

	for m := 0; m < numMethods; m++ {
		methods[m] = instanceValue.Method(m)
	}

	return methods
}

func getTypes(method reflect.Value) []reflect.Type {
	methodArgCount := method.Type().NumIn()
	if methodArgCount == 0 {
		return []reflect.Type{}
	}

	types := make([]reflect.Type, methodArgCount)
	for t := 0; t < methodArgCount; t++ {
		types[t] = method.Type().In(t)
	}

	return types
}
