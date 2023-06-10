package simple_orm

import "context"

type QueryResult struct {
	Result any
	Err    error
}

type QueryContext struct {
}

type MiddleWare func(next HandleFunc) HandleFunc

type HandleFunc func(ctx context.Context, qc *QueryContext) *QueryResult
