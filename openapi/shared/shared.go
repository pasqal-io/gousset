package shared

type Json = any

func Ptr[T any](value T) *T {
	return &value
}

type Nothing = struct{}
