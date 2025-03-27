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
	type SimpleStructNoDescription struct {
		Bool   bool
		Int    int
		Float  float32
		String string
	}
	sample := SimpleStructNoDescription{}
	_, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)
	assert.ErrorContains(t, err, "missing a public name")
}

// Operations should fail if there is no public name.
type SimpleStructWithStructDescriptionButNoPublicName struct {
	Bool   bool
	Int    int
	Float  float32
	String string
}

func (SimpleStructWithStructDescriptionButNoPublicName) Description() string {
	return "This explains what the operation does"
}

var _ doc.HasDescription = SimpleStructWithStructDescriptionButNoPublicName{}

func TestOperationWithStructDescriptionButNoPublicName(t *testing.T) {
	sample := SimpleStructWithStructDescriptionButNoPublicName{}
	_, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)
	assert.ErrorContains(t, err, "missing a public name")
}

// Operations should fail if there is no field summary.

type SimpleStructWithStructDescriptionAndPublicName struct {
	Bool   bool    `path:"bool"`
	Int    int     `path:"int"`
	Float  float32 `path:"float"`
	String string  `path:"string"`
}

func (SimpleStructWithStructDescriptionAndPublicName) Description() string {
	return "This explains what the operation does"
}

var _ doc.HasDescription = SimpleStructWithStructDescriptionAndPublicName{}

func TestOperationWithStructDescriptionAndPublicName(t *testing.T) {
	sample := SimpleStructWithStructDescriptionAndPublicName{}
	_, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)
	assert.ErrorContains(t, err, "doesn't have a description")
}

// Operations should fail if there is no field summary.

type SimpleStructWithDescriptionAndPublicName struct {
	Bool   bool    `path:"bool" summary:"this is a bool" description:"Longer description of a bool"`
	Int    int     `path:"int"  summary:"this is a int" description:"Longer description of a int"`
	Float  float32 `path:"float"  summary:"this is a float" description:"Longer description of a float"`
	String string  `path:"string"  summary:"this is a string" description:"Longer description of a string"`
}

func (SimpleStructWithDescriptionAndPublicName) Description() string {
	return "This explains what the operation does"
}

var _ doc.HasDescription = SimpleStructWithDescriptionAndPublicName{}

func TestOperationWithDescriptionAndPublicName(t *testing.T) {
	sample := SimpleStructWithDescriptionAndPublicName{}
	spec, err := parameter.FromStruct(reflect.TypeOf(sample), parameter.InPath)

	if err != nil {
		t.Fatal(err)
	}

	testutils.EqualJSON(t, spec, `[
  {
    "description": "Longer description of a bool",
    "in": "path",
    "name": "bool",
    "required": true,
    "schema": {
      "type": "boolean"
    }
  },
  {
    "description": "Longer description of a int",
    "in": "path",
    "name": "int",
    "required": true,
    "schema": {
      "format": "int32",
      "type": "number"
    }
  },
  {
    "description": "Longer description of a float",
    "in": "path",
    "name": "float",
    "required": true,
    "schema": {
      "format": "float",
      "type": "number"
    }
  },
  {
    "description": "Longer description of a string",
    "in": "path",
    "name": "string",
    "required": true,
    "schema": {
      "type": "string"
    }
  }
]`)
}
