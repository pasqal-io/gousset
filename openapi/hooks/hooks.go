// Interfaces that users may implement to expand the documentation.
package hooks

import (
	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/schema"
)

type HasSummary = doc.HasSummary
type HasDescription = doc.HasDescription
type HasExternalDocs = doc.HasExternalDocs
type HasExample = example.HasExample
type HasExamples = example.HasExamples
type HasSchema = schema.HasSchema
type HasEnum = schema.HasEnum
type HasFormat = schema.HasEnum
