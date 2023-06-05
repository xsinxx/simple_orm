package simple_orm

type UpsertBuilder[T any] struct {
	insert *Insert[T]
}

type UpsertKey struct {
	assigns []Assignable
}
