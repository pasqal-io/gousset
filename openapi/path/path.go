package path

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/operation"
	"github.com/pasqal-io/gousset/openapi/parameter"
	"github.com/pasqal-io/gousset/openapi/response"
	"github.com/pasqal-io/gousset/openapi/security"
)

// A path in the API.
//
// MUST start with `/`. Path templating is allowed.
type Route string

var routeTemplateRegex = regexp.MustCompile("/:([^/]*)")

func MakeRoute(path string) (Route, error) {
	if !strings.HasPrefix(path, "/") {
		return "<error>", fmt.Errorf("expected a path, starting with '/', got \"%s\"", path)
	}
	// Convert "/:foo"-style captures to "/{foo}"-style captures.
	replaced := routeTemplateRegex.ReplaceAllFunc([]byte(path), func(b []byte) []byte {
		subpath, _ := strings.CutPrefix(string(b), "/:")
		return []byte(fmt.Sprint("/{", strcase.ToSnake(subpath), "}"))
	})
	return Route(replaced), nil
}

// The HTTP verbs.
type Verb string

const (
	Get     = Verb("get")
	Put     = Verb("put")
	Post    = Verb("post")
	Delete  = Verb("delete")
	Options = Verb("options")
	Patch   = Verb("patch")
	Head    = Verb("head")
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
	Parameters *[]parameter.Parameter `json:"parameters,omitempty"`

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
	// The type of inputs.
	//
	// This must be either the zero value (no input) or a struct containing
	// no other field than `Body`, `Query`, `Path`, `Header`.
	Input reflect.Type

	// Security requirements for this endpoint.
	Security []security.Requirement

	// A human-readable summary explaining what this endpoint does.
	Summary string

	// A more detailed description. May contain markdown.
	Description *string

	// Reference to external documentation.
	ExternalDocs *doc.External

	// Information on the response.
	Response response.Implementation

	// If `true`, mark this endpoint as deprecated.
	Deprecated bool
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
			Responses:    verbImpl.Response,
			Deprecated:   verbImpl.Deprecated,
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
