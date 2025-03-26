package testutils

import (
	"encoding/json"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/pasqal-io/gousset/openapi/shared"
)

func EqualJSON(t *testing.T, value shared.Json, reference string) {
	t.Helper()
	ja := jsonassert.New(t)
	valueJson, err := json.Marshal(value)
	if err != nil {
		t.Fatal("Value cannot be marshaled", err)
	}

	ja.Assertf(string(valueJson), reference)
}

type Boolean bool

const (
	BooleanTrue  = Boolean(true)
	BooleanFalse = Boolean(false)
)
