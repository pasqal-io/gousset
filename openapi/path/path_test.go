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
	Numbers []int      `json:"some_numbers" summary:"expecting a few integers"`
	Object  SomeStruct `json:"an_object" summary:"expecting an object"`
}

type MyPath struct {
	Legends []string `path:"some_strings" summary:"expecting a few strings"`
}

type MyQuery struct {
	Numbers []int `path:"path_numbers" query:"query_numbers" summary:"expecting a few integers"`
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
	testutils.EqualJSON(t, result, `{"summary": "Clearly, this is a path",
        "parameters": null,
        "get": {
          "summary": "",
          "operationId": "get /foo/bar",
          "parameters": [
            {
              "name": "some_strings",
              "in": "path",
              "summary": "expecting a few strings",
              "required": true,
              "deprecated": false,
              "schema": {
                "schema": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                }
              }
            },
            {
              "name": "query_numbers",
              "in": "query",
              "summary": "expecting a few integers",
              "required": true,
              "deprecated": false,
              "schema": {
                "schema": {
                  "type": "array",
                  "items": {
                    "type": "number"
                  }
                }
              }
            }
          ]
        },
        "post": {
          "summary": "",
          "operationId": "post /foo/bar",
          "parameters": [
            {
              "name": "some_strings",
              "in": "path",
              "summary": "expecting a few strings",
              "required": true,
              "deprecated": false,
              "schema": {
                "schema": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                }
              }
            },
            {
              "name": "query_numbers",
              "in": "query",
              "summary": "expecting a few integers",
              "required": true,
              "deprecated": false,
              "schema": {
                "schema": {
                  "type": "array",
                  "items": {
                    "type": "number"
                  }
                }
              }
            }
          ],
          "request": {
            "required": true,
            "content": {}
          }
        }
      }
	`)
}
