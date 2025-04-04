package serialization_test

import (
	"testing"

	"github.com/pasqal-io/gousset/inner/serialization"
	"github.com/pasqal-io/gousset/testutils"
)

func TestFlattenStructToJSON(t *testing.T) {
	type InnerStruct struct {
		Text string `json:"text"`
	}
	type OuterStruct struct {
		Field1          int                    `json:"field_1"`
		Field2          bool                   `json:"field_2"`
		NullableField3  *string                `json:"field_3"`
		OmitEmptyField4 *InnerStruct           `json:"field_4,omitempty"`
		InlineMap       map[string]InnerStruct `json:"-1" flatten:""`
		InlineMap2      map[string]InnerStruct `json:"-2" flatten:""`
		OutlineMap      map[string]InnerStruct `json:"map"`
	}
	sample := OuterStruct{
		Field1:          42,
		Field2:          true,
		NullableField3:  nil,
		OmitEmptyField4: nil,
		InlineMap: map[string]InnerStruct{
			"inline_key_1": {
				Text: "text_1",
			},
			"inline_key_2": {
				Text: "text_2",
			},
		},
		InlineMap2: map[string]InnerStruct{
			"inline_key_3": {
				Text: "text_3",
			},
			"inline_key_4": {
				Text: "text_4",
			},
		},
		OutlineMap: map[string]InnerStruct{
			"outline_key_a": {
				Text: "text_a",
			},
			"outline_key_2": {
				Text: "text_b",
			},
		},
	}
	bag, err := serialization.FlattenStructToJSON(sample)
	if err != nil {
		t.Fatal(err)
	}

	testutils.EqualJSON(t, bag, `{
		"field_1": 42,
		"field_2": true,
		"field_3": null,
		"inline_key_1": {
			"text": "text_1"
		},
		"inline_key_2": {
			"text": "text_2"
		},
		"inline_key_3": {
			"text": "text_3"
		},
		"inline_key_4": {
			"text": "text_4"
		},
		"map": {
			"outline_key_2": {
			"text": "text_b"
			},
			"outline_key_a": {
			"text": "text_a"
			}
		}
		}`)
}
