package path

import (
	"github.com/pasqal-io/gousset/openapi/operation"
	"github.com/pasqal-io/gousset/openapi/parameter"
)

// A path in the API.
//
// MUST start with `/`. Path templating is allowed.
type Route = string

type Verb string

const (
	Get     = Verb("get")
	Put     = Verb("put")
	Post    = Verb("post")
	Delete  = Verb("Delete")
	Options = Verb("Options")
)

type Spec struct {
	// Short summary for all operations on this path.
	Summary *string `json:"summary"`

	// Longer description for all operations on this path. May include Markdown.
	Description *string `json:"description"`

	// A list of parameters that are applicable for all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined in the OpenAPI Object’s components.parameters.
	Parameters *[]parameter.Parameter

	PerVerb
}

type PerVerb = map[Verb]operation.Spec
