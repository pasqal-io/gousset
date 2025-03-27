package operation

import (
	"fmt"
	"reflect"

	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/parameter"
	"github.com/pasqal-io/gousset/openapi/request"
	"github.com/pasqal-io/gousset/openapi/response"
	"github.com/pasqal-io/gousset/openapi/security"
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
	SecurityRequirements []security.Requirement `json:"securityRequirements,omitempty"`

	// A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined in the OpenAPI Object’s components.parameters.
	Parameters []parameter.Parameter `json:"parameters,omitempty"`

	// The body expected by this operation.
	Request *request.Request `json:"request,omitempty"`

	// The responses that this operation may return.
	Responses *response.Responses `json:"response,omitempty"`
}

type Implementation struct {
	// The input type.
	//
	// This MUST be a struct.
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
	}

	addParameters := func(field string, in parameter.In, typ reflect.Type) error {
		param, err := parameter.FromStruct(typ, in)
		if err != nil {
			return fmt.Errorf("failed to extract spec for %s of %s: %w", field, operationId, err)
		}
		result.Parameters = append(result.Parameters, param...)
		return nil
	}
	for i := 0; i < impl.Input.NumField(); i++ {
		field := impl.Input.Field(i)
		switch field.Name {
		case "Body":
			request, err := request.FromField(field)
			if err != nil {
				return Spec{}, fmt.Errorf("failed to extract body spec for %s: %w", operationId, err)
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
			return Spec{}, fmt.Errorf("invalid input type %s, it may not have fields other than Path, Query, Header, Body, found %s", impl.Input.String(), field.Name)
		}
	}
	return result, nil
}
