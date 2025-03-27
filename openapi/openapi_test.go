package openapi_test

import (
	"reflect"
	"testing"

	"github.com/pasqal-io/gousset/openapi"
	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/path"
	"github.com/pasqal-io/gousset/openapi/shared"
	"github.com/pasqal-io/gousset/testutils"
	"github.com/pasqal-io/gousset/testutils/structs"
)

// Test on an empty spec.
func TestEmptySpec(t *testing.T) {
	title := "Some title"
	version := "v3.14"
	spec, err := openapi.FromImplementation(openapi.Implementation{
		Info: openapi.Info{
			Title:   title,
			Version: version,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// The string below has been validated with https://editor.swagger.io/, if you change it,
	// please validate the new version.
	testutils.EqualJSONf(t, spec, `
	{
		"openapi": "%s",
		"info": {
			"title": "%s",
			"version": "%s"
		},
		"paths": {},
		"components": {}
	}
  `, openapi.OpenApiVersion, title, version)
}

// Test the Info field.
func TestInfo(t *testing.T) {
	title := "Some title"
	version := "v3.14"
	spec, err := openapi.FromImplementation(openapi.Implementation{
		Info: openapi.Info{
			Title:          title,
			Version:        version,
			Description:    shared.Ptr("This is a description"),
			TermsOfService: shared.Ptr("Very permissive"),
			Contact: &openapi.Contact{
				Name:  "Marvin",
				Url:   shared.Ptr("http://www.example.org/snoopy"),
				Email: shared.Ptr("evil-martian@example.org"),
			},
			License: &openapi.License{
				Name: "We have a custom license",
				Url:  shared.Ptr("http://www.example.org/MPL"),
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	// The string below has been validated with https://editor.swagger.io/, if you change it,
	// please validate the new version.
	testutils.EqualJSONf(t, spec, `{
		"openapi": "%s",
		"info": {
			"title": "Some title",
			"version": "v3.14",
			"description": "This is a description",
			"termsOfService": "Very permissive",
			"contact": {
			"name": "Marvin",
			"url": "http://www.example.org/snoopy",
			"email": "evil-martian@example.org"
			},
			"license": {
			"name": "We have a custom license",
			"url": "http://www.example.org/MPL"
			}
		},
		"paths": {},
		"components": {}
		}`, openapi.OpenApiVersion)
}

// Test external docs

func TestExternalDocs(t *testing.T) {
	spec, err := openapi.FromImplementation(openapi.Implementation{
		ExternalDocs: &doc.External{
			Description: "Look, external docs",
			Url:         "http://www.example.org",
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	// The string below has been validated with https://editor.swagger.io/, if you change it,
	// please validate the new version.
	testutils.EqualJSONf(t, spec, `{
		"openapi": "%s",
		"info": {
			"title": "",
			"version": ""
		},
		"paths": {},
		"components": {},
		"externalDocs": {
			"description": "Look, external docs",
			"url": "http://www.example.org"
		}
		}`, openapi.OpenApiVersion)
}

// Test Endpoints

func TestEndpoints(t *testing.T) {
	type Path struct {
		Foo string `path:"foo" description:"I am foo"`
		Bar string `path:"bar"  description:"I am bar"`
	}
	type Query struct {
		Sna int `query:"sna"  description:"I am sna"`
	}
	type Body struct {
		Ga bool    `json:"ga"  description:"I am ga"`
		Bu []int   `json:"bu"  description:"I am bu"`
		Zo float64 `json:"zo"  description:"I am zo"`
	}
	endPoints := make(map[string]path.Implementation)
	endPoints["/v1/{foo}/{bar}"] = path.Implementation{
		Summary:     "A very interesting endpoint",
		Description: shared.Ptr("With additional description"),
	}
	spec, err := openapi.FromImplementation(openapi.Implementation{
		Endpoints: []path.Implementation{
			{
				Path:        "/v1/{foo}/{bar}",
				Summary:     "A very interesting endpoint",
				Description: shared.Ptr("With additional description"),
				PerVerb: map[path.Verb]path.VerbImplementation{
					path.Get: {
						Input: reflect.TypeOf(structs.PathQuery[Path, Query]{
							Path:  Path{},
							Query: Query{},
						}),
						Summary: "This is the summary for GET /foo/bar",
					},
					path.Put: {
						Input: reflect.TypeOf(structs.BodyPathQuery[Body, Path, Query]{
							Path:  Path{},
							Query: Query{},
							Body:  Body{},
						}),
						Summary: "This is the summary for PUT /foo/bar",
					},
				},
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	// The string below has been validated with https://editor.swagger.io/, if you change it,
	// please validate the new version.
	//testutils.ValidateOpenAPI(t, spec)
	testutils.EqualJSONf(t, spec, `{
		"openapi": "%s",
		"info": {
			"title": "",
			"version": ""
		},
		"paths": {
			"/v1/{foo}/{bar}": {
			"summary": "A very interesting endpoint",
			"description": "With additional description",
			"get": {
				"summary": "This is the summary for GET /foo/bar",
				"operationId": "get /v1/{foo}/{bar}",
				"parameters": [
				{
					"description": "I am foo",
					"in": "path",
					"name": "foo",
					"required": true,
					"schema": {
					"type": "string"
					}
				},
				{
					"description": "I am bar",
					"in": "path",
					"name": "bar",
					"required": true,
					"schema": {
					"type": "string"
					}
				},
				{
					"description": "I am sna",
					"in": "query",
					"name": "sna",
					"required": true,
					"schema": {
					"format": "int32",
					"type": "number"
					}
				}
				],
				"responses": {
				"default": {
					"description": ""
				}
				}
			},
			"put": {
				"summary": "This is the summary for PUT /foo/bar",
				"operationId": "put /v1/{foo}/{bar}",
				"parameters": [
				{
					"description": "I am foo",
					"in": "path",
					"name": "foo",
					"required": true,
					"schema": {
					"type": "string"
					}
				},
				{
					"description": "I am bar",
					"in": "path",
					"name": "bar",
					"required": true,
					"schema": {
					"type": "string"
					}
				},
				{
					"description": "I am sna",
					"in": "query",
					"name": "sna",
					"required": true,
					"schema": {
					"format": "int32",
					"type": "number"
					}
				}
				],
				"requestBody": {
				"required": true,
				"content": {
					"application/json": {
					"schema": {
						"type": "object",
						"required": [
						"ga",
						"bu",
						"zo"
						],
						"properties": {
						"bu": {
							"type": "array",
							"items": {
							"type": "number",
							"format": "int32"
							}
						},
						"ga": {
							"type": "boolean"
						},
						"zo": {
							"type": "number",
							"format": "double"
						}
						}
					}
					}
				}
				},
				"responses": {
				"default": {
					"description": ""
				}
				}
			}
			}
		},
		"components": {}
		}`, openapi.OpenApiVersion)
}
