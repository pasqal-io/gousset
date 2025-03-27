package openapi

import (
	"errors"
	"fmt"

	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/header"
	"github.com/pasqal-io/gousset/openapi/link"
	"github.com/pasqal-io/gousset/openapi/parameter"
	"github.com/pasqal-io/gousset/openapi/path"
	"github.com/pasqal-io/gousset/openapi/response"
	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/openapi/security"
)

type RequestBody struct {
	Description string

	// Determine if the request body is required in the request.
	Required bool
}

// Contact information for the exposed API.
type Contact struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Email string `json:"email"`
}

type License struct {
	// The license name used for the API.
	Name string `json:"name"`

	// An [SPDX-Licenses] expression for the API. The identifier field is mutually exclusive of the url field.
	Identifier *string `json:"identifier"`

	// A URI for the license used for the API. This MUST be in the form of a URI. The url field is mutually exclusive of the identifier field.
	Url *string `json:"url"`
}

// General information on the API.
type Info struct {
	// The title of the API.
	Title string `json:"title"`

	// Version of the API.
	Version string `json:"version"`

	// Short summary of the API.
	Summary *string `json:"summary"`

	// Longer description. May include Markdown.
	Description *string `json:"description"`

	// A URI for the Terms of Service for the API. This MUST be in the form of a URI.
	TermsOfService *string `json:"termsOfService"`

	// Contact information for the exposed API.
	Contact *Contact `json:"contact"`

	// License information for the exposed API.
	License *License `json:"license"`
}

// A specification, fit for consumption by an OpenAPI client.
type Spec struct {
	// The version number of the OpenAPI spec.
	OpenApiVersion string `json:"openapi"`

	// General information on the API.
	Info Info `json:"info"`

	// All the routes covered by this API.
	Paths map[path.Route]path.Spec `json:"paths"`

	// A set of reusable objects for different aspects of the OAS. All objects defined within the Components Object will have no effect on the API unless they are explicitly referenced from outside the Components Object
	Components Components `json:"components"`

	// Additional external documentations.
	ExternalDocs *doc.External `json:"externalDocs"`
}

type Components struct {
	Schemas         *map[string]schema.Schema
	Responses       *map[string]response.Spec
	Parameters      *map[string]parameter.Parameter
	Examples        *map[string]example.Example
	SecuritySchemes *map[string]security.Scheme
	Headers         *map[string]header.Header
	Links           *map[string]link.Link
}

// User-provided metadata containing information on the implementation
// to be converted to OpenAPI spec.
type Implementation struct {
	// General information on the service.
	Info Info

	// Endpoints within this service.
	Endpoints map[string]path.Implementation

	// Any additional external documentation.
	ExternalDocs *doc.External `json:"externalDocs"`
}

// Build a complete OpenAPI spec from a description of an implementation.
func FromImplementation(implem Implementation) (Spec, error) {
	result := Spec{
		OpenApiVersion: "3.1.1",
		Info:           implem.Info,
		ExternalDocs:   implem.ExternalDocs,
	}
	paths := make(map[path.Route]path.Spec)
	result.Paths = paths
	for onePath, pathImpl := range implem.Endpoints {
		route, err := path.MakeRoute(onePath)
		if err != nil {
			return Spec{}, errors.New("invalid path")
		}
		pathSpec, err := path.FromPath(pathImpl)
		if err != nil {
			return Spec{}, fmt.Errorf("failed to build spec for route %s: %w", route, err)
		}
		paths[route] = pathSpec
	}
	return result, nil
}
