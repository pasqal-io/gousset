package example

import "github.com/pasqal-io/gousset/openapi/shared"

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

type HasExample interface {
	Example() shared.Json
}

type HasExamples interface {
	Examples() Example
}
