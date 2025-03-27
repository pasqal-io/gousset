package structs

type Body[B any] struct {
	Body B
}

type Query[Q any] struct {
	Query Q
}

type Path[P any] struct {
	Path P
}

type BodyQuery[B any, Q any] struct {
	Body  B
	Query Q
}

type BodyPath[B any, P any] struct {
	Body B
	Path P
}

type PathQuery[P any, Q any] struct {
	Path  P
	Query Q
}

type BodyPathQuery[B any, P any, Q any] struct {
	Body  B
	Path  P
	Query Q
}
