package link

type Link interface {
	sealed()
}

// https://spec.openapis.org/oas/v3.0.1.html#link-object
type Spec struct {
	OperationRef string
	// ... Other fields to be added as needed
}

func (Spec) sealed() {}

var _ Link = Spec{}

type Reference struct {
	Ref string `json:"$ref"`
}

func Ref(to string) Reference {
	return Reference{
		Ref: to,
	}
}

func (Reference) sealed() {}

// User-provided metadata containing information on the implementation
// to be converted to OpenAPI spec.
type Implementation struct {
	OperationRef string
}

func FromImplementation(impl Implementation) (Link, error) {
	result := Spec{}
	result.OperationRef = impl.OperationRef
	return result, nil
}
