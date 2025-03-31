package shared

type Json = any

// Take a pointer to any value.
func Ptr[T any](value T) *T {
	return &value
}
