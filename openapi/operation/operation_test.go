package operation_test

import (
	"reflect"
	"testing"

	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/operation"
	"github.com/pasqal-io/gousset/openapi/parameter"
	"github.com/pasqal-io/gousset/openapi/security"
	"github.com/pasqal-io/gousset/testutils"
	"gotest.tools/assert"
)

// Operations should fail if there is no summary.
func TestOperationWithoutSummary(t *testing.T) {
	type SimpleStructNoSummary struct {
		Bool   bool
		Int    int
		Float  float32
		String string
	}
	sample := SimpleStructNoSummary{}
	_, err := operation.FromStruct(reflect.TypeOf(sample), []security.Requirement{}, parameter.InPath, "GET", "/foo/bar")
	assert.ErrorContains(t, err, "HasSummary")
}

// Operations should fail if there is no public name.

type SimpleStructWithStructSummaryButNoPublicName struct {
	Bool   bool
	Int    int
	Float  float32
	String string
}

func (SimpleStructWithStructSummaryButNoPublicName) Summary() string {
	return "This explains what the operation does"
}

var _ doc.HasSummary = SimpleStructWithStructSummaryButNoPublicName{}

func TestOperationWithStructSummaryButNoPublicName(t *testing.T) {
	sample := SimpleStructWithStructSummaryButNoPublicName{}
	_, err := operation.FromStruct(reflect.TypeOf(sample), []security.Requirement{}, parameter.InPath, "GET", "/foo/bar")
	assert.ErrorContains(t, err, "missing a public name")
}

// Operations should fail if there is no field summary.

type SimpleStructWithStructSummaryAndPublicName struct {
	Bool   bool    `path:"bool"`
	Int    int     `path:"int"`
	Float  float32 `path:"float"`
	String string  `path:"string"`
}

func (SimpleStructWithStructSummaryAndPublicName) Summary() string {
	return "This explains what the operation does"
}

var _ doc.HasSummary = SimpleStructWithStructSummaryAndPublicName{}

func TestOperationWithStructSummaryAndPublicName(t *testing.T) {
	sample := SimpleStructWithStructSummaryAndPublicName{}
	_, err := operation.FromStruct(reflect.TypeOf(sample), []security.Requirement{}, parameter.InPath, "GET", "/foo/bar")
	assert.ErrorContains(t, err, "doesn't have a summary")
}

// Operations should succeed if we have everything!

type SimpleStructWithSummaryAndPublicName struct {
	Bool   bool    `path:"bool" summary:"this is a bool"`
	Int    int     `path:"int"  summary:"this is a int"`
	Float  float32 `path:"float"  summary:"this is a float"`
	String string  `path:"string"  summary:"this is a string"`
}

func (SimpleStructWithSummaryAndPublicName) Summary() string {
	return "This explains what the operation does"
}

var _ doc.HasSummary = SimpleStructWithSummaryAndPublicName{}

func TestOperationWithSummaryAndPublicName(t *testing.T) {
	sample := SimpleStructWithSummaryAndPublicName{}
	spec, err := operation.FromStruct(reflect.TypeOf(sample), []security.Requirement{}, parameter.InPath, "GET", "/foo/bar")

	if err != nil {
		t.Fatal(err)
	}

	testutils.EqualJSON(t, spec, `{
		"operationId": "GET /foo/bar path",
		"securityRequirements": [],
		"parameters": [{
			"deprecated": false,
			"description": null,
			"in": "path",
			"name": "bool",
			"required": true,
			"schema": {
			  "schema": {
				"type":"bool"
			  }
			},
			"summary": "this is a bool"
        }, {
			"deprecated": false,
			"description": null,
			"in": "path",
			"name": "int",
			"required": true,
			"schema": {
			  "schema": {
				"type":"number"
			  }
			},
			"summary": "this is a int"
        }, {
			"deprecated": false,
			"description": null,
			"in": "path",
			"name": "float",
			"required": true,
			"schema": {
			  "schema": {
				"type":"number"
			  }
			},
			"summary": "this is a float"
        }, {
			"deprecated": false,
			"description": null,
			"in": "path",
			"name": "string",
			"required": true,
			"schema": {
			  "schema": {
				"type":"string"
			  }
			},
			"summary": "this is a string"
        }],
		"summary": "This explains what the operation does"
	}`)
}
