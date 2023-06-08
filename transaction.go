package simple_orm

import (
	"context"
	"database/sql"
	"github.com/simple_orm/model"
	"github.com/simple_orm/valuer"
)

// db & tx 均对这个接口做了实现
type session interface {
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type TX struct {
	core
	tx *sql.Tx
	db *DB
}

type core struct {
	r       *model.Registry
	creator valuer.Creator // 运行时再执行的函数，默认unsafe
	dialect Dialect        // 方言
}

func (t *TX) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args)
}

func (t *TX) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args)
}

func (t *TX) getCore() core {
	return t.db.core
}

func (t *TX) Commit() error {
	return t.tx.Commit()
}

func (t *TX) Rollback() error {
	return t.tx.Rollback()
}
