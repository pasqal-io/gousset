package media

import (
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/openapi/shared"
)

type Type struct {
	Schema   *schema.Schema     `json:"schema"`
	Example  *shared.Json       `json:"example"`
	Examples *[]example.Example `json:"examples"`
}
