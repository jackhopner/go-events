package events

import (
	"reflect"
	"testing"
)

type TestInterface interface {
	Event_One(evt *testType)
	Two()
}

type testType struct{}
type testInterface struct{}

func (t *testInterface) Event_One(evt *testType) {}
func (t *testInterface) Two()                    {}

func newTestInterface() TestInterface {
	return &testInterface{}
}

func Test_getMethods(t *testing.T) {
	o := newTestInterface()
	instanceType := reflect.TypeOf(o)

	methods := getMethods(o)

	numMethods := len(methods)
	if numMethods != 2 {
		t.Fatalf("Expected to get 2 methods but got %d", numMethods)
	}

	method0Str := methods[0].Type().String()
	expected0Str := "func(*events.testType)"
	if method0Str != expected0Str {
		t.Fatalf(
			"Expected string to be [%s] was [%s]",
			expected0Str,
			method0Str,
		)
	}

	method0Name := instanceType.Method(0).Name
	expected0Name := "Event_One"
	if method0Name != expected0Name {
		t.Fatalf(
			"Expected name to be [%s] was [%s]",
			expected0Name,
			method0Name,
		)
	}

	method1Str := methods[1].Type().String()
	expected1Str := "func()"
	if method1Str != expected1Str {
		t.Fatalf(
			"Expected string to be [%s] was [%s]",
			expected1Str,
			method1Str,
		)
	}

	method1Name := instanceType.Method(1).Name
	expected1Name := "Two"
	if method1Name != expected1Name {
		t.Fatalf(
			"Expected name to be [%s] was [%s]",
			expected1Name,
			method1Name,
		)
	}
}

func Test_getTypes(t *testing.T) {
	o := newTestInterface()
	// instanceValue := reflect.ValueOf(o)

	methods := getMethods(o)

	method0Types := getTypes(methods[0])

	num0Types := len(method0Types)
	expectedNum0Types := 1
	if num0Types != expectedNum0Types {
		t.Fatalf(
			"Expected %d types but got %d",
			expectedNum0Types,
			num0Types,
		)
	}

	expected0Type := reflect.TypeOf(&testType{})
	method0Type := method0Types[0]
	if method0Type != expected0Type {
		t.Fatalf("Expected type [%+v] but got [%+v]", expected0Type, method0Type)
	}

	num1Types := len(getTypes(methods[1]))
	expectedNum1Types := 0
	if num1Types != expectedNum1Types {
		t.Fatalf(
			"Expected %d types but got %d",
			expectedNum1Types,
			num1Types,
		)
	}
}
