package simple_orm

import "reflect"

type DBOption func(db *DB)

type DB struct {
	r *Registry
}

func NewDB(opts ...DBOption) (*DB, error) {
	db := &DB{
		r: &Registry{
			models: map[reflect.Type]*tableModel{},
		},
	}
	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}

func MustNewDB(opts ...DBOption) *DB {
	db, err := NewDB(opts...)
	// init panic
	if err != nil {
		panic(err)
	}
	return db
}

func DBWithRegister(r *Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func RegisterWithModels(models map[reflect.Type]*tableModel) DBOption {
	return func(db *DB) {
		db.r.models = models
	}
}
