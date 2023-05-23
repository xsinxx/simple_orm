package simple_orm

import (
	"database/sql"
	"reflect"
)

type DBOption func(db *DB)

type DB struct {
	r     *Registry
	store *sql.DB // 对应具体数据库的存储
}

func OpenDB(store *sql.DB, opts ...DBOption) (*DB, error) {
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
