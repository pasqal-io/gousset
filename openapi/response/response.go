package response

import (
	"encoding/json"
	"fmt"

	"github.com/pasqal-io/gousset/inner/serialization"
	"github.com/pasqal-io/gousset/openapi/header"
	"github.com/pasqal-io/gousset/openapi/link"
	"github.com/pasqal-io/gousset/openapi/media"
	"github.com/pasqal-io/gousset/openapi/shared"
)

type Response interface {
	sealed()
}

type Reference struct {
	Ref string `json:"$ref"`
}

func Ref(to string) Reference {
	return Reference{
		Ref: to,
	}
}
func (Reference) sealed() {}

var _ Response = Ref("")

type Responses struct {
	Default Response             `json:"default"`
	PerCode *map[uint16]Response `json:"-,omitempty" flatten:""`
}

func (r Responses) MarshalJSON() ([]byte, error) {
	var nilResponse Response = nil
	if r.Default == nilResponse {
		return []byte{}, fmt.Errorf("error while flattening Responses for serialization: Default cannot be nil")
	}
	flattened, err := serialization.FlattenStructToJSON(r)
	if err != nil {
		return []byte{}, fmt.Errorf("error while flattening Responses for serialization: %w, ", err)
	}
	return json.Marshal(flattened)
}

var _ json.Marshaler = Responses{}

type Spec struct {
	Description string                    `json:"description"`
	Headers     *map[string]header.Header `json:"headers,omitempty"`
	Content     *map[string]media.Type    `json:"content,omitempty"`
	Links       *map[string]link.Link     `json:"links,omitempty"`
}

func (Spec) sealed() {}

var _ Response = Spec{}

// User-provided metadata containing information on the implementation
// to be converted to OpenAPI spec.
type Implementation struct {
	Default ResponseImplementation
	PerCode *map[uint16]ResponseImplementation
}

func FromImplementation(impl Implementation) (Responses, error) {
	def, err := FromResponseImplementation(impl.Default)
	if err != nil {
		return Responses{}, fmt.Errorf("while compiling response, error in default response: %w", err)
	}
	result := Responses{
		Default: def,
	}
	if impl.PerCode != nil {
		perCode := make(map[uint16]Response)
		for k, v := range *impl.PerCode {
			response, err := FromResponseImplementation(v)
			if err != nil {
				return Responses{}, fmt.Errorf("while compiling response, error in response code %d: %w", k, err)
			}
			perCode[k] = response
		}
		result.PerCode = &perCode
	}
	return result, nil
}

type ResponseImplementation struct {
	Description string
	Headers     *map[string]header.Implementation
	Content     *map[string]media.Implementation
	Links       *map[string]link.Implementation
}

func FromResponseImplementation(impl ResponseImplementation) (Response, error) {
	result := Spec{
		Description: impl.Description,
	}
	if impl.Headers != nil {
		result.Headers = shared.Ptr(make(map[string]header.Header))
		for k, v := range *impl.Headers {
			h, err := header.FromImplementation(v)
			if err != nil {
				return result, fmt.Errorf("while compiling header %s, error: %w", k, err)
			}
			(*result.Headers)[k] = h
		}
	}
	if impl.Content != nil {
		result.Content = shared.Ptr(make(map[string]media.Type))
		for k, v := range *impl.Content {
			h, err := media.FromImplementation(v)
			if err != nil {
				return result, fmt.Errorf("while compiling media type %s, error: %w", k, err)
			}
			(*result.Content)[k] = h
		}
	}
	if impl.Links != nil {
		result.Links = shared.Ptr(make(map[string]link.Link))
		for k, v := range *impl.Links {
			h, err := link.FromImplementation(v)
			if err != nil {
				return result, fmt.Errorf("while compiling link %s, error: %w", k, err)
			}
			(*result.Links)[k] = h
		}
	}
	return result, nil
}
