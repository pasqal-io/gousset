package request

import (
	"fmt"
	"reflect"

	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/media"
	"github.com/pasqal-io/gousset/openapi/schema"
)

type Request interface {
	sealed()
}

type Spec struct {
	// A brief description of the request body. This could contain examples of use. [CommonMark] syntax MAY be used for rich text representation.
	Description *string `json:"description,omitempty"`
	Required    bool    `json:"required"`

	// The content of the request body. The key is a media type or media type range, see [RFC7231] Appendix D, and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. "text/plain" overrides "text/*"
	Content map[string]media.Type `json:"content"`
}

func (Spec) sealed() {}

// Reference to a body described in the OpenAPI components.
type Reference struct {
	Ref         string  `json:"$ref"`
	Summary     *string `json:"summary,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (Reference) sealed() {}

var _ Request = Reference{}

func FromField(from reflect.StructField) (Request, error) {
	// Extract summary and description.
	description := doc.GetDescription(from.Type)
	content := make(map[string]media.Type)
	schema, err := schema.FromImplementation(schema.Implementation{Type: from.Type, PublicNameKey: "json"})
	if err != nil {
		return Spec{}, fmt.Errorf("failed to extract the type of body %s: %w", from.Type.String(), err)
	}

	// FIXME: Support example.

	content["application/json"] = media.Type{
		Schema: &schema,
	}
	return Spec{
		Description: description,
		Required:    from.Type.Kind() != reflect.Ptr,
		// We do not support content yet.
		Content: content,
	}, nil
}
