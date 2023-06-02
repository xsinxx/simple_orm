package simple_orm

type OnDuplicateKeyBuilder[T any] struct {
	insert *Insert[T]
}

type OnDuplicateKey struct {
	assigns []Assignable
}
