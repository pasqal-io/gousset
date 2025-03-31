// Examples provided as part of the OpenAPI spec.
package example

import (
	"reflect"

	"github.com/pasqal-io/gousset/openapi/shared"
)

// Example Object | Reference Object
type Example interface {
	sealed()
}

// https://spec.openapis.org/oas/v3.0.1.html#example-object
type Spec struct {
	Summary       string
	Description   string
	Value         *shared.Json
	ExternalValue *string
}

func (Spec) sealed() {}

var _ Example = Spec{}

// A reference to a Component.
type Reference shared.Reference

func Ref(to string) Reference {
	return Reference(shared.Ref(to))
}

func (Reference) sealed() {}

var _ Example = Reference{}

// Implement this to add a single example of a type without comments.
type HasExample interface {
	Example() shared.Json
}

// Imlpement this to add examples to a type with comments.
type HasExamples interface {
	Examples() map[string]Example
}

// Utility cast to HasExample, get the example.
func GetExample(typ reflect.Type) *shared.Json {
	phony := reflect.New(typ)
	if !phony.CanInterface() {
		return nil
	}
	if casted, ok := phony.Interface().(HasExample); ok {
		return shared.Ptr(casted.Example())
	}
	return nil
}

// Utility cast to HasExamples, get the examples.
func GetExamples(typ reflect.Type) *map[string]Example {
	phony := reflect.New(typ)
	if !phony.CanInterface() {
		return nil
	}
	if casted, ok := phony.Interface().(HasExamples); ok {
		return shared.Ptr(casted.Examples())
	}
	return nil
}
