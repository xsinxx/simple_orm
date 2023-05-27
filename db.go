package simple_orm

import (
	"database/sql"
	"github.com/simple_orm/model"
	"github.com/simple_orm/valuer"
	"reflect"
)

type DBOption func(db *DB)

type DB struct {
	store   *sql.DB // 对应具体数据库的存储, 继承
	r       *model.Registry
	creator valuer.Creator // 运行时再执行的函数，默认unsafe
}

func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(store *sql.DB, opts ...DBOption) (*DB, error) {
	db := &DB{
		store: store,
		r: &model.Registry{
			TableModels: map[reflect.Type]*model.TableModel{},
		},
		creator: valuer.NewUnsafeValue,
	}
	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}

func DBWithRegister(r *model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

// DBWithCreator 提供了unsafe & reflect 两种方式将sql中读取的数据放入到结构体指针中
func DBWithCreator(reflectCreator valuer.Creator) DBOption {
	return func(db *DB) {
		db.creator = reflectCreator
	}
}
