package testutils

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/pasqal-io/gousset/openapi/shared"
	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
)

func ValidateOpenAPI(t *testing.T, value shared.Json) {
	valueJson, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic("invalid json")
	}
	document, err := libopenapi.NewDocument(valueJson)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to create document: %w", err))
	}

	docValidator, validatorErrs := validator.NewValidator(document)
	if validatorErrs != nil {
		for _, err := range validatorErrs {
			t.Log(fmt.Errorf("invalid OpenAPI document: %w", err))
		}
		t.Fatal("document failed validation stage 2")
	}

	valid, validationErrs := docValidator.ValidateDocument()

	if !valid {
		for _, e := range validationErrs {
			// 5. Handle the error
			fmt.Printf("Type: %s, Failure: %s\n", e.ValidationType, e.Message)
			fmt.Printf("Fix: %s\n\n", e.HowToFix)
		}
		t.Fatal("document failed validation stage 3")
	}
}

func EqualJSON(t *testing.T, value shared.Json, reference string) {
	t.Helper()
	ja := jsonassert.New(t)
	valueJson, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic("invalid json")
	}
	// We dislpay the JSON, even if it takes a few milliseconds.
	//
	// `go test` will swallow this from the logs
	fmt.Println("Checking json", string(valueJson))
	ja.Assert(string(valueJson), reference)
}

func EqualJSONf(t *testing.T, value shared.Json, reference string, others ...any) {
	t.Helper()
	ja := jsonassert.New(t)
	valueJson, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(fmt.Errorf("invalid json (candidate value): %w", err))
	}
	// We dislpay the JSON, even if it takes a few milliseconds.
	//
	// `go test` will swallow this from the logs
	fmt.Println("Checking json", string(valueJson))
	ja.Assertf(string(valueJson), reference, others...)
}

func PrintJSON(value shared.Json) {
	valueJson, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic("invalid json")
	}
	fmt.Println("Checking json", string(valueJson))
}

type Boolean bool

const (
	BooleanTrue  = Boolean(true)
	BooleanFalse = Boolean(false)
)
