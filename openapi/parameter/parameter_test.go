package parameter_test

import (
	"reflect"
	"testing"

	"github.com/pasqal-io/gousset/openapi/parameter"
	"github.com/pasqal-io/gousset/testutils"
	"gotest.tools/assert"
)

// Operations should succeed if there is no public name, with a default public name.
//
// Note: Warnings expected.
func TestParameterMissingRename(t *testing.T) {
	type SimpleStruct struct {
		Bool   bool
		Int    int
		Float  float32
		String string
	}
	sample := SimpleStruct{}
	parameters, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)
	assert.NilError(t, err)

	testutils.EqualJSON(t, parameters, `[
		{
			"in": "path",
			"name": "bool",
			"required": true,
			"schema": {
			"type": "boolean"
			}
		},
		{
			"in": "path",
			"name": "int",
			"required": true,
			"schema": {
			"format": "int32",
			"type": "number"
			}
		},
		{
			"in": "path",
			"name": "float",
			"required": true,
			"schema": {
			"format": "float",
			"type": "number"
			}
		},
		{
			"in": "path",
			"name": "string",
			"required": true,
			"schema": {
			"type": "string"
			}
		}
		]`)
}

// Operations should succeed if there is a public name and adopt this public name.
//
// Note: Warnings expected.
type SimpleStructWithPublicName struct {
	Bool   bool    `path:"my_bool"`
	Int    int     `path:"my_int"`
	Float  float32 `path:"my_float"`
	String string  `path:"my_string"`
}

func TestParameterWithStructDescriptionAndPublicName(t *testing.T) {
	sample := SimpleStructWithPublicName{}
	parameters, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)
	assert.NilError(t, err)

	testutils.EqualJSON(t, parameters, `[
		{
			"in": "path",
			"name": "my_bool",
			"required": true,
			"schema": {
			"type": "boolean"
			}
		},
		{
			"in": "path",
			"name": "my_int",
			"required": true,
			"schema": {
			"format": "int32",
			"type": "number"
			}
		},
		{
			"in": "path",
			"name": "my_float",
			"required": true,
			"schema": {
			"format": "float",
			"type": "number"
			}
		},
		{
			"in": "path",
			"name": "my_string",
			"required": true,
			"schema": {
			"type": "string"
			}
		}
		]`)
}

// Operations should succeed if there is a description and adopt this description.
//
// No warnings expected.
type SimpleStructWithDescriptionAndPublicName struct {
	Bool   bool    `path:"my_bool" description:"Longer description of a bool"`
	Int    int     `path:"my_int"  description:"Longer description of a int"`
	Float  float32 `path:"my_float"  description:"Longer description of a float"`
	String string  `path:"my_string"  description:"Longer description of a string"`
}

func TestParameterWithDescriptionAndPublicName(t *testing.T) {
	sample := SimpleStructWithDescriptionAndPublicName{}
	spec, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)

	if err != nil {
		t.Fatal(err)
	}

	testutils.EqualJSON(t, spec, `[
		{
			"description": "Longer description of a bool",
			"in": "path",
			"name": "my_bool",
			"required": true,
			"schema": {
			"type": "boolean"
			}
		},
		{
			"description": "Longer description of a int",
			"in": "path",
			"name": "my_int",
			"required": true,
			"schema": {
			"format": "int32",
			"type": "number"
			}
		},
		{
			"description": "Longer description of a float",
			"in": "path",
			"name": "my_float",
			"required": true,
			"schema": {
			"format": "float",
			"type": "number"
			}
		},
		{
			"description": "Longer description of a string",
			"in": "path",
			"name": "my_string",
			"required": true,
			"schema": {
			"type": "string"
			}
		}
		]`)
}
