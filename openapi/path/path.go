package path

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/operation"
	"github.com/pasqal-io/gousset/openapi/parameter"
)

// A path in the API.
//
// MUST start with `/`. Path templating is allowed.
type Route string

func MakeRoute(path string) (Route, error) {
	if strings.HasPrefix(path, "/") {
		return Route(path), nil
	}
	return "<error>", fmt.Errorf("expected a path, starting with '/', got \"%s\"", path)
}

// The HTTP verbs.
type Verb string

const (
	Get     = Verb("get")
	Put     = Verb("put")
	Post    = Verb("post")
	Delete  = Verb("delete")
	Options = Verb("options")
)

// Specifications for one path (all verbs).
type Spec struct {
	// Short summary for all operations on this path.
	Summary string `json:"summary,omitempty"`

	// Longer description for all operations on this path. May include Markdown.
	Description *string `json:"description,omitempty"`

	// A list of parameters that are applicable for all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined in the OpenAPI Object’s components.parameters.
	//
	// In the current implementation, we expect that this contains the path parameters.
	Parameters *[]parameter.Parameter `json:"parameters"`

	Get     *operation.Spec `json:"get,omitempty"`
	Put     *operation.Spec `json:"put,omitempty"`
	Post    *operation.Spec `json:"post,omitempty"`
	Delete  *operation.Spec `json:"delete,omitempty"`
	Options *operation.Spec `json:"options,omitempty"`
}

// User-provided metadata containing information on the implementation
// to be converted to OpenAPI spec (all verbs at one path).
type Implementation struct {
	Summary     string
	Description *string
	Path        string
	PerVerb     map[Verb]VerbImplementation
}

// User-provided metadata containing information on the implementation
// to be converted to OpenAPI spec (a single verb at one path).
type VerbImplementation struct {
	Input        reflect.Type
	Security     []security.Requirement
	Summary      string
	Description  *string
	ExternalDocs *doc.External
}

func FromPath(impl Implementation) (Spec, error) {
	result := Spec{
		Summary:     impl.Summary,
		Description: impl.Description,
		// We do not attempt to factorize shared parameters.
		Parameters: nil,
	}
	for verb, verbImpl := range impl.PerVerb {
		op, err := operation.FromImplementation(operation.Implementation{
			Input:        verbImpl.Input,
			Verb:         string(verb),
			Path:         impl.Path,
			Security:     verbImpl.Security,
			Summary:      verbImpl.Summary,
			Description:  verbImpl.Description,
			ExternalDocs: verbImpl.ExternalDocs,
		})
		if err != nil {
			return Spec{}, fmt.Errorf("failed to extract specs for operation %s at %s: %w", verb, impl.Path, err)
		}
		var ptr **operation.Spec
		switch verb {
		case Get:
			ptr = &result.Get
		case Put:
			ptr = &result.Put
		case Post:
			ptr = &result.Post
		case Delete:
			ptr = &result.Delete
		case Options:
			ptr = &result.Options
		default:
			panic(fmt.Sprint("Verb not handled ", verb))
		}
		*ptr = &op
	}
	return result, nil
}
