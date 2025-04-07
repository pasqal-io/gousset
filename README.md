# Gousset

Gousset, the Go Usable Semi-static Schema Extractor Tool, is a tool designed to extract OpenAPI specs from Go data structures.

## Usage

To generate an OpenAPI spec, call `openapi.FromImplementation` as follows:

```go
import 	"github.com/pasqal-io/gousset/openapi"
import 	"github.com/pasqal-io/gousset/openapi/operation"
import 	"github.com/pasqal-io/gousset/openapi/path"
import 	"github.com/pasqal-io/gousset/openapi/response"
import 	"github.com/pasqal-io/gousset/shared"

spec, err := openapi.FromImplementation(openapi.Implementation {
    Info: openapi.Info {
        // Fill in the metadata for your API, explaining what it's for,
        // who to contact for assistance, etc.
        // ...
    },
    Endpoints: []path.Implementation {
        {
            Path: "/api/v1/user/{id}",
            Summary: "Operations on an individual user",
            PerVerb: map[Verb]operation.Impleementation {
                path.Get: {
                    Summary: "Get the information on the user",
                    Response: response.Implementation {
                        Default: response.ResponseImplementation {
                            Description: "In case of success",
                            Contents: &map[string]media.Implementation {
                                "application/json": {
                                    Example: &MyUser {
                                        Name: "John Doe",
                                    },
                                },
                            },
                        },
                    },
                }
            },
        }
    },
})
```

This is a bit heavy, but if your code or framework is sufficiently high-level, you should be
able to extract the information automatically from the code.

## Conventions

### Renaming

Use tags `json`, `query`, `path`, `header` to rename go-style UpperCamelCased fields into
their corresponding public names.

### Flattening

Use tag `flatten` to flatten a struct or a map into its container, e.g.

```go
type User struct {
    Name string `json:"name"`
    Email string `json:"email"`
}

type MyArguments struct {
    User User `flatten:""`
}
```

to document this as if this were the following struct

```go
type MyArguments struct {
    Name string `json:"name"`
    Email string `json:"email"`
}
```

This is useful when you share `User` between several structs. This works also through maps
and interfaces.

## Documenting parameters/responses.

gousset recognizes the following tags that you may use to document individual parameters:

### Description (recommended)

Use tag `description` to explain the parameter. This may use markdown.

```go
type MyArguments struct {
    Arg Type `description:"some description"`
}
```

### Format (recommended whenever possible)

Use tag `format` on a string or number to restrain the parameter to some format, e.g.

```go
type MyArguments struct {
    Email string `format:"email"`
}
```

See also `HasFormat` for an alternative method if you use a `type Email string` instead of
a raw `string`.

See `gousset.openapi.schema.Format*` constants for a list of weel-known formats.

### Deprecation

Use tag `deprecated` on a field to mark it as deprecated.

### Example (recommended)

Use tag `example` on a string to provide an example of its value.

See also `HasExample` and `HasExamples` to generalize this to an entire type.

### Pattern

Use tag `pattern` to restrict a string field to some regex.

### Min/max

Use tags `minimum:"number"`, `maximum:"number"`, `exclusiveMinimum:"number"`,
`exclusiveMaximum:"number"` to restrict the values of numbers.

### Length
Use `maxItems:"number"`, `minItems:"number"` to restrict the length of arrays.

Use `maxLength:"number"`, `minLength:"number"` to restrict the length of strings.

Use `maxProperties:"number"`, `minProperties:"number"` to restrict the number of
    properties of a map.


## Customizing types

You can also customize entire types.

### Sum types

Sum types (e.g. `int | string`) are very different between Go and OpenSpec.

In Go, se represent the sum type as a `struct`, e.g.

```go
type MyFooOrBar struct {
    Foo* Foo `variant:"variant_1" json:"foo"`
    Bar* Bar `variant:"variant_2" json:"bar"`
}
```

This will be compiled as if we had defined

```ts
type MyFooOrBar = {Foo* Foo} | {Bar* Bar}
```

Alternatively, you can request flattening

```go
type MyFooOrBar struct {
    Foo* Foo `variant:"variant_1" json:"foo" flatten:""`
    Bar* Bar `variant:"variant_2" json:"bar"  flatten:""`
}
```

This will be compiled as if we had defined

```ts
type MyFooOrBar = Foo | Bar
```

As of this writing, this mechanism only works if `Foo` or `Bar` are `struct`. If you
need it to work with other types, don't hesitate to file an issue!

### Min, max, pattern, length, ...

See all the interfaces in `hooks` to see how to document entire types.
