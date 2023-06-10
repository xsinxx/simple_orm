package simple_orm

import (
	"context"
	"database/sql"
	"github.com/hashicorp/go-multierror"
	"github.com/simple_orm/model"
	"github.com/simple_orm/valuer"
	"reflect"
)

type DBOption func(db *DB)

type DB struct {
	core          // 元数据信息
	store *sql.DB // 对应具体数据库的存储
}

type TxKey struct {
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
		core: core{
			r: &model.Registry{
				TableModels: map[reflect.Type]*model.TableModel{},
			},
			creator: valuer.NewUnsafeValue,
			dialect: mySQLDialect,
		},
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

// DBWithDialect 默认是MySQL，支持配置
func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}

// 开启事务应归属于db，事务的提交回滚属于事务
func (db *DB) beginTx(ctx context.Context, opts *sql.TxOptions) (*TX, error) {
	tx, err := db.store.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &TX{
		tx: tx,
		db: db,
	}, nil
}

// 事务扩散，当ctx中无事务则直接开启需求
func (db *DB) beginTxIfNotExists(ctx context.Context, opts *sql.TxOptions) (context.Context, *TX, error) {
	txInContext := ctx.Value(TxKey{})
	if txInContext != nil {
		t := txInContext.(*TX)
		return ctx, t, nil
	}
	tx, err := db.store.BeginTx(ctx, opts)
	if err != nil {
		return ctx, nil, err
	}
	ctx = context.WithValue(ctx, TxKey{}, tx)
	return ctx, &TX{tx: tx, db: db}, nil
}

// 事务闭包：当执行事务出错或执行中发生panic需要回滚
func (db *DB) doTx(ctx context.Context, task func(ctx2 context.Context, tx *TX) error, opts *sql.TxOptions) (err error) {
	tx, err := db.beginTx(ctx, opts)
	if err != nil {
		return err
	}
	panicked := true
	defer func() {
		if !panicked && err == nil {
			return
		}
		// 回滚中出现问题则将error追加
		txErr := tx.Rollback()
		if txErr != nil {
			err = multierror.Append(err, txErr)
		}
	}()
	// 执行任务
	err = task(ctx, tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	panicked = false
	return nil
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.store.ExecContext(ctx, query, args...)
}

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.store.QueryContext(ctx, query, args...)
}

func (db *DB) getCore() core {
	return db.core
}
