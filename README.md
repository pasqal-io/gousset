# Goussset

Goussset, the Go Usable Semi-Static Schema Extractor Tool, is a tool designed to extract OpenAPI specs from Go data structures.

## Usage

Implement `HasSummary`, `HasDocumentation`, `HasExternalDocs`, `HasExample`, `HasExamples` to add documentation to your data structures.

Add a tag `summary` to your fields to add documentation to an individual field.

Implement `HasSchema` to implement more complex cases, e.g. enums.