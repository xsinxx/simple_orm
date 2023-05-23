package simple_orm

import (
	"database/sql"
	"github.com/simple_orm/valuer"
	"reflect"
)

type DBOption func(db *DB)

type DB struct {
	*sql.DB       // 对应具体数据库的存储, 继承
	r             *Registry
	unsafeCreator valuer.Creator // 运行时再执行的函数
}

func OpenDB(store *sql.DB, opts ...DBOption) (*DB, error) {
	db := &DB{
		r: &Registry{
			models: map[reflect.Type]*TableModel{},
		},
		unsafeCreator: valuer.NewUnsafeValue,
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

func DBWithUnsafeCreator(unsafeCreator valuer.Creator) DBOption {
	return func(db *DB) {
		db.unsafeCreator = unsafeCreator
	}
}
