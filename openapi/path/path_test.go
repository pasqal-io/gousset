package path_test

import (
	"reflect"
	"testing"

	"github.com/pasqal-io/gousset/openapi/path"
	"github.com/pasqal-io/gousset/testutils"
	"gotest.tools/assert"
)

// Test that Route is always prefixed with a "/".
func TestMakeRoute(t *testing.T) {
	_, err := path.MakeRoute("foo/bar")
	assert.ErrorContains(t, err, "expected a path")

	route, err := path.MakeRoute("/foo/bar")
	assert.NilError(t, err)

	assert.Equal(t, string(route), "/foo/bar")
}

// Test that we can compile the OpenAPI Spec from a fairly simple implementation.

type SomeStruct struct {
	Int    int    `json:"an_int"`
	String string `json:"a_string"`
}
type MyBody struct {
	Numbers []int      `json:"some_numbers" description:"expecting a few integers"`
	Object  SomeStruct `json:"an_object" description:"expecting an object"`
}

type MyPath struct {
	Legends []string `path:"some_strings" description:"expecting a few strings"`
}

type MyQuery struct {
	Numbers []int `path:"path_numbers" query:"query_numbers" description:"expecting a few integers"`
}

type BodyPathQuery struct {
	Body  MyBody
	Path  MyPath
	Query MyQuery
}

type PathQuery struct {
	Path  MyPath
	Query MyQuery
}

func TestFromPath(t *testing.T) {
	perVerb := make(map[path.Verb]path.VerbImplementation)
	perVerb[path.Get] = path.VerbImplementation{
		Input: reflect.TypeOf(PathQuery{}),
	}
	perVerb[path.Post] = path.VerbImplementation{
		Input: reflect.TypeOf(BodyPathQuery{}),
	}
	result, err := path.FromPath(path.Implementation{
		Summary: "Clearly, this is a path",
		Path:    "/foo/bar",
		PerVerb: perVerb,
	})
	if err != nil {
		t.Fatal(err)
	}
	testutils.EqualJSON(t, result, `{
		"summary": "Clearly, this is a path",
		"get": {
			"summary": "",
			"operationId": "get /foo/bar",
			"parameters": [
			{
				"description": "expecting a few strings",
				"in": "path",
				"name": "some_strings",
				"required": true,
				"schema": {
				"items": {
					"type": "string"
				},
				"type": "array"
				}
			},
			{
				"description": "expecting a few integers",
				"in": "query",
				"name": "query_numbers",
				"required": true,
				"schema": {
				"items": {
					"type": "number",
					"format": "int32"
				},
				"type": "array"
				}
			}
			],
			"responses": {
			"default": {
				"description": ""
			}
			}
		},
		"post": {
			"summary": "",
			"operationId": "post /foo/bar",
			"parameters": [
			{
				"description": "expecting a few strings",
				"in": "path",
				"name": "some_strings",
				"required": true,
				"schema": {
				"items": {
					"type": "string"
				},
				"type": "array"
				}
			},
			{
				"description": "expecting a few integers",
				"in": "query",
				"name": "query_numbers",
				"required": true,
				"schema": {
				"items": {
					"type": "number",
					"format": "int32"
				},
				"type": "array"
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
					"some_numbers",
					"an_object"
					],
					"properties": {
					"an_object": {
						"type": "object",
						"required": [
						"an_int",
						"a_string"
						],
						"properties": {
						"a_string": {
							"type": "string"
						},
						"an_int": {
							"type": "number",
							"format": "int32"
						}
						}
					},
					"some_numbers": {
						"type": "array",
						"items": {
						"type": "number",
						"format": "int32"
						}
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
	}`)
}
