// Interfaces that users may implement to expand the documentation.
package hooks

import (
	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/schema"
)

// Implement this to add a summary to all instances of a type.
type HasSummary = doc.HasSummary

// Implement this to add a description to all instances of a type.
type HasDescription = doc.HasDescription

// Implement this to add external docs to all instances of a type.
type HasExternalDocs = doc.HasExternalDocs

// Implement this to add an example without comments to all instances of a type.
type HasExample = example.HasExample

// Implement this to add a examples with comments to all instances of a type.
type HasExamples = example.HasExamples

// Implement this to override the schema that gousset infers for a type.
type HasSchema = schema.HasSchema

// Implement this on enums of constants to let gousset list all the possibilities.
type HasEnum = schema.IsEnum

// Implement this on numbers or strings to restrict to a specific format.
type HasFormat = schema.HasFormat

// Implement this to mark a minimal length for a string.
type HasMinStringLength = schema.HasMinStringLength

// Implement this to mark a maximal length for a string.
type HasMaxStringLength = schema.HasMaxStringLength

// Implement this to mark a minimal length for an array.
type HasMinArrayLength = schema.HasMinArrayLength

// Implement this to mark a maximal length for an array.
type HasMaxArrayLength = schema.HasMaxArrayLength

// Implement this to mark a minimal length for an array.
type HasMinMapLength = schema.HasMinMapLength

// Implement this to mark a maximal length for an array.
type HasMaxMapLength = schema.HasMaxMapLength

// Implement this to mark a minimal value for a number.
type HasMin = schema.HasMin

// Implement this to mark a maximal value for a number.
type HasMax = schema.HasMax
