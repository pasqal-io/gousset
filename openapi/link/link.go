package link

type Link interface {
	sealed()
}

type Spec struct {
	OperationRef string
	// ... TBD
}

func (Spec) sealed() {}

var _ Link = Spec{}

type Reference string

func Ref(to string) Reference {
	return Reference(to)
}

func (Reference) sealed() {}
