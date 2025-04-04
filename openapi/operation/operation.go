package operation

import (
	"fmt"
	"reflect"

	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/parameter"
	"github.com/pasqal-io/gousset/openapi/request"
	"github.com/pasqal-io/gousset/openapi/response"
	"github.com/pasqal-io/gousset/openapi/security"
	"github.com/pasqal-io/gousset/shared/structs"
)

type Spec struct {
	// Short summary.
	Summary string `json:"summary"`

	// Longer description. May include Markdown.
	Description *string `json:"description,omitempty"`

	// Additional external documentations.
	ExternalDocs *doc.External `json:"externalDocs,omitempty"`

	// Unique id, case sensitive.
	OperationId string `json:"operationId"`

	// A declaration of which security mechanisms can be used for this operation. The list of values includes alternative Security Requirement Objects that can be used. Only one of the Security Requirement Objects need to be satisfied to authorize a request. To make security optional, an empty security requirement ({}) can be included in the array. This definition overrides any declared top-level security. To remove a top-level security declaration, an empty array can be used.
	SecurityRequirements []security.Requirement `json:"security,omitempty"`

	// A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined in the OpenAPI Object’s components.parameters.
	Parameters []parameter.Parameter `json:"parameters,omitempty"`

	// The body expected by this operation.
	Request *request.Request `json:"requestBody,omitempty"`

	// The responses that this operation may return.
	Responses response.Responses `json:"responses"`

	// If true, this endpoint is deprecated.
	Deprecated bool `json:"deprecated,omitempty"`
}

// User-provided metadata containing information on the implementation
// to be converted to OpenAPI spec.
type Implementation struct {
	// The input type.
	//
	// This MUST be a struct or zero.
	//
	// It MUST NOT contain fields other than
	// - Path
	// - Query
	// - Body
	// - Header
	Input        reflect.Type
	Verb         string
	Path         string
	Security     []security.Requirement
	Summary      string
	Description  *string
	ExternalDocs *doc.External
	Responses    response.Implementation
	Deprecated   bool
}

// Extract an OpenAPI spec for an operation from a description of the implementation.
func FromImplementation(impl Implementation) (Spec, error) {
	operationId := fmt.Sprint(impl.Verb, " ", impl.Path)

	result := Spec{
		Summary:              impl.Summary,
		Description:          impl.Description,
		ExternalDocs:         impl.ExternalDocs,
		OperationId:          operationId,
		SecurityRequirements: impl.Security,
		Deprecated:           impl.Deprecated,
	}

	addParameters := func(field string, in parameter.In, typ reflect.Type) error {
		param, err := parameter.FromStruct(typ, in)
		if err != nil {
			return fmt.Errorf("while compiling operation %s %s, failed to extract spec for %s of %s: %w", impl.Verb, impl.Path, field, operationId, err)
		}
		result.Parameters = append(result.Parameters, param...)
		return nil
	}
	// Zero value, assume the empty struct.
	if impl.Input.Kind() == reflect.Invalid {
		// Note: `impl` is passed by copy, so this mutation is not observable.
		impl.Input = reflect.TypeOf(structs.Nothing{})
	}
	for i := 0; i < impl.Input.NumField(); i++ {
		field := impl.Input.Field(i)
		switch field.Name {
		case "Body":
			request, err := request.FromField(field)
			if err != nil {
				return Spec{}, fmt.Errorf("while compiling operation %s %s, failed to extract body spec for %s: %w", impl.Verb, impl.Path, operationId, err)
			}
			result.Request = &request
		case "Path":
			err := addParameters(field.Name, parameter.InPath, field.Type)
			if err != nil {
				return Spec{}, err
			}
		case "Query":
			err := addParameters(field.Name, parameter.InQuery, field.Type)
			if err != nil {
				return Spec{}, err
			}
		case "Header":
			err := addParameters(field.Name, parameter.InHeader, field.Type)
			if err != nil {
				return Spec{}, err
			}
		default:
			return Spec{}, fmt.Errorf("while compiling operation %s %s, invalid input type %s, it may not have fields other than Path, Query, Header, Body, found %s", impl.Verb, impl.Path, impl.Input.String(), field.Name)
		}
	}
	responses, err := response.FromImplementation(impl.Responses)
	if err != nil {
		return Spec{}, fmt.Errorf("while compiling operation %s %s, invalid response: %w", impl.Verb, impl.Path, err)
	}
	result.Responses = responses
	return result, nil
}
