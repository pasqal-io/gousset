package example

import (
	"reflect"

	"github.com/pasqal-io/gousset/openapi/shared"
)

type Example interface {
	sealed()
}
type Spec struct {
	Summary       string
	Description   string
	Value         *shared.Json
	ExternalValue *string
}

func (Spec) sealed() {}

var _ Example = Spec{}

type Reference struct {
	Ref         string `json:"$ref"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

// A field/struct/... that has a single example without further comments.
type HasExample interface {
	Example() shared.Json
}

// A field/struct/... that has detailed examples.
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
