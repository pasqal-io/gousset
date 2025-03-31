package parameter_test

import (
	"reflect"
	"testing"

	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/parameter"
	"github.com/pasqal-io/gousset/testutils"
	"gotest.tools/assert"
)

// Operations should fail if there is no public name.
func TestOperationMissingRawName(t *testing.T) {
	type SimpleStructNoSummary struct {
		Bool   bool
		Int    int
		Float  float32
		String string
	}
	sample := SimpleStructNoSummary{}
	_, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)
	assert.ErrorContains(t, err, "missing a public name")
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
	_, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)
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
	_, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)
	assert.ErrorContains(t, err, "doesn't have a summary")
}

// Operations should succeed if everything is present.

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
	spec, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)

	if err != nil {
		t.Fatal(err)
	}

	testutils.EqualJSON(t, spec, `[{
			"deprecated": false,
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
			"in": "path",
			"name": "string",
			"required": true,
			"schema": {
			  "schema": {
				"type":"string"
			  }
			},
			"summary": "this is a string"
        }]`)
}
