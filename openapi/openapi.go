// Extract OpenAPI specifications from an existing API.
//
// Use function `openapi.FromImplementation` to convert
// an existing API to an OpenAPI specification. You may
// cache this specification and expose it as an entrypoint
// to generate a user-friendly documentation.
package openapi

import (
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

// The version of OpenAPI specs we're based on.
const OpenApiVersion = "3.0.1"

// Contact information for the exposed API.
type Contact struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name"`

	// The URI for the contact information. This MUST be in the form of a URI.
	Url *string `json:"url,omitempty"`

	// The email address of the contact person/organization. This MUST be in the form of an email address.
	Email *string `json:"email,omitempty"`
}

// License information for the exposed API.
type License struct {
	// The license name used for the API.
	Name string `json:"name"`

	// An [SPDX-Licenses] expression for the API. The identifier field is mutually exclusive of the url field.
	Identifier *string `json:"identifier,omitempty"`

	// A URI for the license used for the API. This MUST be in the form of a URI. The url field is mutually exclusive of the identifier field.
	Url *string `json:"url,omitempty"`
}

// General information on the API.
type Info struct {
	// The title of the API.
	Title string `json:"title"`

	// Version of the API.
	Version string `json:"version"`

	// Longer description. May include Markdown.
	Description *string `json:"description,omitempty"`

	// A URI for the Terms of Service for the API. This MUST be in the form of a URI.
	TermsOfService *string `json:"termsOfService,omitempty"`

	// Contact information for the exposed API.
	Contact *Contact `json:"contact,omitempty"`

	// License information for the exposed API.
	License *License `json:"license,omitempty"`
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
	ExternalDocs *doc.External `json:"externalDocs,omitempty"`
}

// A library of objects that may be used throughout the spec.
//
// Each of these tables is a mapping id -> schema. Refer to them
// with a Reference.
type Components struct {
	// Schema objects.
	Schemas *map[string]schema.Schema `json:"schemas,omitempty"`

	// Response objects.
	Responses *map[string]response.Spec `json:"responses,omitempty"`

	// Parameter objects.
	Parameters *map[string]parameter.Parameter `json:"parameters,omitempty"`

	// Example objects.
	Examples *map[string]example.Example `json:"examples,omitempty"`

	// Security scheme objects.
	SecuritySchemes *map[string]security.Scheme `json:"securitySchemes,omitempty"`

	// Header objects.
	Headers *map[string]header.Header `json:"headers,omitempty"`

	// Link objects.
	Links *map[string]link.Link `json:"links,omitempty"`
}

// User-provided metadata containing information on the implementation
// to be converted to OpenAPI spec.
type Implementation struct {
	// General information on the service.
	Info Info

	// Endpoints within this service.
	Endpoints []path.Implementation

	// Any additional external documentation.
	ExternalDocs *doc.External
}

// Build a complete OpenAPI spec from a description of an implementation.
func FromImplementation(implem Implementation) (Spec, error) {
	result := Spec{
		OpenApiVersion: OpenApiVersion,
		Info:           implem.Info,
		ExternalDocs:   implem.ExternalDocs,
	}
	paths := make(map[path.Route]path.Spec)
	result.Paths = paths
	for _, pathImpl := range implem.Endpoints {
		route, err := path.MakeRoute(pathImpl.Path)
		if err != nil {
			return Spec{}, fmt.Errorf("invalid path: %w", err)
		}
		pathSpec, err := path.FromPath(pathImpl)
		if err != nil {
			return Spec{}, fmt.Errorf("failed to build spec for route %s: %w", route, err)
		}
		paths[route] = pathSpec
	}
	return result, nil
}
