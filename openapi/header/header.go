package header

import (
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/media"
	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/openapi/shared"
)

type Header interface {
	sealed()
}

type Reference string

func Ref(to string) Reference {
	return Reference(to)
}
func (Reference) sealed() {}

var _ Header = Reference("")

type Spec struct {
	Description *string `json:"description"`
	Required    bool    `json:"required"`
	Deprecated  bool    `json:"deprecated"`
	*SchemaSpec
	*ContentSpec
}

type SchemaSpec struct {
	// Describes how the parameter value will be serialized depending on the type of the parameter value. Default values (based on value of in): for "query" - "form"; for "path" - "simple"; for "header" - "simple"; for "cookie" - "form".
	Style string `json:"style"`

	// When this is true, parameter values of type array or object generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this field has no effect. When style is "form", the default value is true. For all other styles, the default value is false. Note that despite false being the default for deepObject, the combination of false with deepObject is undefined.
	Explode bool `json:"explode"`

	Schema schema.Schema `json:"schema"`

	Example  *shared.Json       `json:"example"`
	Examples *[]example.Example `json:"examples"`
}

type ContentSpec struct {
	Content map[string]media.Type `json:"content"`
}
