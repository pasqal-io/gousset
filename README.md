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

To add 

Implement `HasSummary`, `HasDocumentation`, `HasExternalDocs`, `HasExample`, `HasExamples` to add documentation to your data structures.

Add a tag `summary` to your fields to add documentation to an individual field.

Implement `HasSchema` to implement more complex cases, e.g. enums.