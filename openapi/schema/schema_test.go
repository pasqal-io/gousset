package schema_test

import (
	"reflect"
	"testing"

	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/openapi/shared"
	"github.com/pasqal-io/gousset/testutils"
)

// Check the schema for a simple boolean.
func TestBool(t *testing.T) {
	sample := false
	result, err := schema.FromType(reflect.TypeOf(sample), "json")
	if err != nil {
		t.Fatal(err)
	}
	testutils.EqualJSON(t, result, `{
		"type": "bool"
	}`)
}

// Check the schema for a private type implemented by boolean.
func TestCustomBool(t *testing.T) {
	type MyBool bool
	sample := MyBool(false)
	result, err := schema.FromType(reflect.TypeOf(sample), "json")
	if err != nil {
		t.Fatal(err)
	}
	testutils.EqualJSON(t, result, `{
		"type": "bool"
	}`)
}

// Check the schema for a private type implemented by boolean.
func TestCustomBool2(t *testing.T) {
	sample := testutils.BooleanFalse
	result, err := schema.FromType(reflect.TypeOf(sample), "json")
	if err != nil {
		t.Fatal(err)
	}
	testutils.EqualJSON(t, result, `{
		"type": "bool"
	}`)
}

// Check that HasSchema is used.
type BoolWithEmptySchema bool

func (BoolWithEmptySchema) Schema() schema.Schema {
	return schema.OneOf{OneOf: []schema.Schema{}}
}

var _ schema.HasSchema = BoolWithEmptySchema(true)

func TestBoolWithHasSchema(t *testing.T) {
	sample := BoolWithEmptySchema(false)
	result, err := schema.FromType(reflect.TypeOf(sample), "json")
	if err != nil {
		t.Fatal(err)
	}
	testutils.EqualJSON(t, result, `{
	    "oneOf": []
	}`)
}

// Check that HasExample is used.
type BoolWithExample bool

func (BoolWithExample) Example() shared.Json {
	return true
}

var _ example.HasExample = BoolWithExample(true)

func TestBoolWithHasExample(t *testing.T) {
	sample := BoolWithExample(false)
	result, err := schema.FromType(reflect.TypeOf(sample), "json")
	if err != nil {
		t.Fatal(err)
	}
	testutils.EqualJSON(t, result, `{
		"example": true,
		"type": "bool"
	}`)
}

// Check that HasExternalDocs is used.
type BoolWithExternalDocs bool

func (BoolWithExternalDocs) Docs() doc.External {
	return doc.External{
		Description: "Look, there's something interesting over there",
		Url:         "http://www.example.org",
	}
}

var _ doc.HasExternalDocs = BoolWithExternalDocs(true)

func TestBoolWithExternalDocs(t *testing.T) {
	sample := BoolWithExternalDocs(false)
	result, err := schema.FromType(reflect.TypeOf(sample), "json")
	if err != nil {
		t.Fatal(err)
	}
	testutils.EqualJSON(t, result, `{
	    "externalDocs": {
		    "description": "Look, there's something interesting over there",
			"url": "http://www.example.org"
		},
		"type": "bool"
	}`)
}
