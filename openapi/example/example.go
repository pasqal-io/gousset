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

// A field/struct/... that has a single example without further comments.
type HasExample interface {
	Example() shared.Json
}

// A field/struct/... that has detailed examples.
type HasExamples interface {
	Examples() Example
}
