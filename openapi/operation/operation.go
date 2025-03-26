package operation

import (
	"fmt"
	"reflect"

	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/parameter"
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
	SecurityRequirements []security.Requirement `json:"securityRequirements"`

	// A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined in the OpenAPI Object’s components.parameters.
	Parameters []parameter.Parameter `json:"parameters"`
}

func FromStruct(Struct reflect.Type, security []security.Requirement, in parameter.In, verb string, path string) (Spec, error) {
	if Struct.Kind() != reflect.Struct {
		return Spec{}, fmt.Errorf("invalid type %s, expected a struct, got %v.", Struct.String(), Struct.Kind())
	}

	summary := doc.GetSummary(Struct)
	if summary == nil {
		return Spec{}, fmt.Errorf("invalid type %s, should implement HasSummary", Struct.String())
	}
	description := doc.GetDescription(Struct)
	externalDocs := doc.GetExternalDocs(Struct)

	var parameters []parameter.Parameter
	for i := 0; i < Struct.NumField(); i++ { // We have checked above that it's a struct.
		field := Struct.Field(i)
		// FIXME: We'll need to know if there are any default values.
		param, err := parameter.FromField(field, in)
		if err != nil {
			return Spec{}, fmt.Errorf("failed to generate spec for parameter %s of type %s: %w", field.Name, Struct.String(), err)
		}
		parameters = append(parameters, param)

	}

	return Spec{
		Summary:              *summary, // Non-nil, checked above.
		Description:          description,
		ExternalDocs:         externalDocs,
		OperationId:          fmt.Sprint(verb, " ", path, " ", in),
		Parameters:           parameters,
		SecurityRequirements: security,
	}, nil
}
