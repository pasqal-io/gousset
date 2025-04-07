package schema_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/shared"
)

// -----   Example: Using sum types.

type Unflattened[T any, U any] struct {
	Comments string `json:"comments"`
	Ok       *T     `json:"ok" variant:"ok"`
	Error    *U     `json:"error" variant:"error"`
}

func ExampleSchema() {
	spec, err := schema.FromImplementation(schema.Implementation{
		Type:          reflect.TypeFor[Unflattened[int, bool]](),
		PublicNameKey: "json",
	})
	if err != nil {
		panic(err)
	}
	marshaled, err := json.Marshal(spec)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(marshaled))
	// Output: {"oneOf":[{"type":"object","required":["comments","ok"],"properties":{"comments":{"type":"string"},"ok":{"type":"number","format":"int32"}}},{"type":"object","required":["comments","error"],"properties":{"comments":{"type":"string"},"error":{"type":"boolean"}}}]}
}

type Flattened[T any, U any] struct {
	Comments string `json:"comments"`
	Ok       *T     `flatten:"" variant:"ok"`
	Error    *U     `flatten:"" variant:"error"`
}

func TestFlattened(t *testing.T) {
	type Success struct {
		Success string `json:"success"`
	}
	type Failure struct {
		Failure bool `json:"failure"`
	}
	spec, err := schema.FromImplementation(schema.Implementation{
		Type:          reflect.TypeFor[Flattened[Success, Failure]](),
		PublicNameKey: "json",
	})
	if err != nil {
		panic(err)
	}
	// Note: we cannot test with testutils.EqualJSON because the
	// order of fields in `Required` changes across runs.
	err = schema.EqualSchema(spec, schema.OneOf{
		OneOf: []schema.Schema{
			schema.Object{
				Shared: schema.Shared{
					Type: schema.TypeObject,
				},
				Required: []string{
					"comments",
					"success",
				},
				Properties: map[string]schema.Schema{
					"comments": schema.Primitive{
						Shared: schema.Shared{
							Type: schema.TypeString,
						},
					},
					"success": schema.Primitive{
						Shared: schema.Shared{
							Type: schema.TypeString,
						},
					},
				},
			},
			schema.Object{
				Shared: schema.Shared{
					Type: schema.TypeObject,
				},
				Required: []string{
					"comments",
					"failure",
				},
				Properties: map[string]schema.Schema{
					"comments": schema.Primitive{
						Shared: schema.Shared{
							Type: schema.TypeString,
						},
					},
					"failure": schema.Primitive{
						Shared: schema.Shared{
							Type: schema.TypeBool,
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMoreFlattened(t *testing.T) {
	type Success struct {
		Success string `json:"result"`
	}
	type Failure struct {
		Failure bool `json:"result"`
	}
	spec, err := schema.FromImplementation(schema.Implementation{
		Type:          reflect.TypeFor[Flattened[Success, Failure]](),
		PublicNameKey: "json",
	})
	if err != nil {
		panic(err)
	}

	// Note: we cannot test with testutils.EqualJSON because the
	// order of fields in `Required` changes across runs.
	err = schema.EqualSchema(spec, schema.OneOf{
		OneOf: []schema.Schema{
			schema.Object{
				Shared: schema.Shared{
					Type: schema.TypeObject,
				},
				Required: []string{
					"comments",
					"result",
				},
				Properties: map[string]schema.Schema{
					"comments": schema.Primitive{
						Shared: schema.Shared{
							Type: schema.TypeString,
						},
					},
					"result": schema.Primitive{
						Shared: schema.Shared{
							Type: schema.TypeString,
						},
					},
				},
			},
			schema.Object{
				Shared: schema.Shared{
					Type: schema.TypeObject,
				},
				Required: []string{
					"comments",
					"result",
				},
				Properties: map[string]schema.Schema{
					"comments": schema.Primitive{
						Shared: schema.Shared{
							Type: schema.TypeString,
						},
					},
					"result": schema.Primitive{
						Shared: schema.Shared{
							Type: schema.TypeBool,
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFullyFlattened(t *testing.T) {
	type Int struct {
		Field int `json:"field"`
	}
	type String struct {
		Field string `json:"field"`
	}
	type Sum struct {
		*Int    `variant:"int" flatten:""`
		*String `variant:"string" flatten:""`
	}
	spec, err := schema.FromImplementation(schema.Implementation{
		Type:          reflect.TypeFor[Sum](),
		PublicNameKey: "json",
	})
	if err != nil {
		panic(err)
	}
	// Note: we cannot test with testutils.EqualJSON because the
	// order of fields in `Required` changes across runs.
	err = schema.EqualSchema(spec, schema.OneOf{
		OneOf: []schema.Schema{
			schema.Object{
				Shared: schema.Shared{
					Type: schema.TypeObject,
				},
				Required: []string{"field"},
				Properties: map[string]schema.Schema{
					"field": schema.Primitive{
						Shared: schema.Shared{
							Type:   schema.TypeNumber,
							Format: shared.Ptr(string(schema.FormatInt32)),
						},
					},
				},
			},
			schema.Object{
				Shared: schema.Shared{
					Type: schema.TypeObject,
				},
				Required: []string{"field"},
				Properties: map[string]schema.Schema{
					"field": schema.Primitive{
						Shared: schema.Shared{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
