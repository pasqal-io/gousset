package doc

import (
	"reflect"

	"github.com/pasqal-io/gousset/openapi/shared"
)

type External struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

type HasSummary interface {
	Summary() string
}

type HasDescription interface {
	Description() string
}

type HasExternalDocs interface {
	Docs() External
}

// Utility: get the summary from a HasSummary.
func GetSummary(typ reflect.Type) *string {
	phony := reflect.New(typ)
	if !phony.CanInterface() {
		return nil
	}
	if hasSummary, ok := phony.Interface().(HasSummary); ok {
		return shared.Ptr(hasSummary.Summary())
	}
	return nil
}

// Utility: get the description from a HasDescription.
func GetDescription(typ reflect.Type) *string {
	phony := reflect.New(typ)
	if !phony.CanInterface() {
		return nil
	}
	if hasDescription, ok := phony.Interface().(HasDescription); ok {
		return shared.Ptr(hasDescription.Description())
	}
	return nil
}

// Utility: get the description from a HasDescription.
func GetExternalDocs(typ reflect.Type) *External {
	phony := reflect.New(typ)
	if !phony.CanInterface() {
		return nil
	}
	if hasExternalDocs, ok := phony.Interface().(HasExternalDocs); ok {
		return shared.Ptr(hasExternalDocs.Docs())
	}
	return nil
}
