package simple_orm

import (
	"context"
)

type RawQuery[T any] struct {
	Builder         // builder是Insert & Select公共部分
	core            // core中是元数据信息
	session session // session是db或tx
	sql     string  // 原生查询的SQL
	args    []any   // 原生查询的参数
}

func NewRawQuery[T any](session session, sql string, args ...any) *RawQuery[T] {
	return &RawQuery[T]{
		core:    session.getCore(),
		session: session,
		sql:     sql,
		args:    args,
	}
}

func (r *RawQuery[T]) Build() (*Query, error) {
	return &Query{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}

func (r *RawQuery[T]) Get(ctx context.Context) (*T, error) {
	return get[T](ctx, r.core, r.session, &QueryContext{
		Builder: r,
	})
}

func (r *RawQuery[T]) GetMul(ctx context.Context) ([]*T, error) {
	return getMul[T](ctx, r.core, r.session, &QueryContext{
		Builder: r,
	})
}
