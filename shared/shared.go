package shared

type Json = any

// Take a pointer to any value.
func Ptr[T any](value T) *T {
	return &value
}

// A reference to a definition provided within Components.
type Reference struct {
	Ref string `json:"$ref"`
}

func Ref(to string) Reference {
	return Reference{
		Ref: to,
	}
}
