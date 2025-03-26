package response

import (
	"github.com/pasqal-io/gousset/openapi/header"
	"github.com/pasqal-io/gousset/openapi/link"
	"github.com/pasqal-io/gousset/openapi/media"
)

type Response interface {
	sealed()
}

type Reference string

func Ref(to string) Reference {
	return Reference(to)
}
func (Reference) sealed() {}

var _ Response = Reference("")

type Responses struct {
	Default Response `json:"default"`
	*PerCode
}

type PerCode = map[uint16]Response

type Spec struct {
	Description string                    `json:"description"`
	Headers     *map[string]header.Header `json:"headers"`
	Content     *map[string]media.Type    `json:"content"`
	Links       *map[string]link.Link     `json:"links"`
}

func (Spec) sealed() {}

var _ Response = Spec{}
