// Serialization utilities.
package serialization

import (
	"fmt"
	"reflect"

	"github.com/pasqal-io/gousset/inner/tags"
)

// Flatten a struct to a JSON object.
//
// # We have many structs that look like
//
//	type Foo struct {
//	   Field1 A
//	   Field2 B
//	   SomeMap map[string]Foobar
//	}
//
// # That are meant to be converted to JSON as
//
//	{
//	   "field1": ...,
//	   "field2": ...,
//	   "key1": ..., // Key of SomeMap
//	   "key2": ..., // Key of SomeMap
//	   ...
//	}
//
// This function handles conversion to a JSON object that
// can then be serialized using the regular `json.Marshal`.
//
// Conventions:
//
//   - all field names MUST have a tag `json`.
//   - a field holding a map MAY have a tag `flatten`, in which case
//     the map is flattened into the containing struct, as above
//   - a field hodling a struct MAY have a tag `flatten`, in which case
//     the struct is flattened into the containing struct
//   - a field MAY specify `omitempty` as the second value of `json`
//     to specify that zero values should be skipped
//
// This function should never panic.
func FlattenStructToJSON(value any) (map[string]any, error) {
	result := make(map[string]any)
	reflected := reflect.ValueOf(value)
	if reflected.Type().Kind() != reflect.Struct {
		return result, fmt.Errorf("while flattening to json, invalid type, expected a struct, got %s", reflected.Type().String())
	}
	for i := 0; i < reflected.Type().NumField(); i++ {
		field := reflected.Field(i)
		fieldTags, err := tags.Parse(reflected.Type().Field(i).Tag)
		if err != nil {
			return result, fmt.Errorf("while flattening to json, couldn't parse tags of field %s of type %s: %w", field.Type().Name(), reflected.Type().String(), err)
		}
		if field.IsZero() && fieldTags.ShouldOmitEmpty("json") {
			continue
		}

		if fieldTags.IsFlattened() {
			for {
				if field.IsZero() && fieldTags.ShouldOmitEmpty("json") {
					// Handle `nil`.
					break
				}
				switch field.Type().Kind() {
				case reflect.Map:
					iter := field.MapRange()
					for iter.Next() {
						result[iter.Key().String()] = iter.Value().Interface()
					}
				case reflect.Struct:
					flattened, err := FlattenStructToJSON(field.Interface())
					if err != nil {
						return result, err
					}
					for k, v := range flattened {
						result[k] = v
					}
				case reflect.Interface:
					fallthrough
				case reflect.Pointer:
					field = field.Elem()
					continue
				default:
					return result, fmt.Errorf("while flattening to json, in type %s, found a field marked as `flatten`, expecting a struct or a map, got %s", reflected.Type().String(), field.Type().String())
				}
				break
			}
		} else {
			name := fieldTags.PublicFieldName("json")
			if name == nil {
				return result, fmt.Errorf("while flattening to json, field %s of type %s is missing a `json` tag", field.Type().Name(), reflected.Type().String())
			}

			if !field.IsZero() || !fieldTags.ShouldOmitEmpty("json") {
				result[*name] = field.Interface()
			}
		}
	}

	return result, nil
}
