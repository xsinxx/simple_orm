package simple_orm

import (
	"context"
	"errors"
)

func get[T any](ctx context.Context, core core, session session, qc *QueryContext) (*T, error) {
	var handler HandleFunc = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getHandler[T](ctx, session, core, qc)
	}
	middlewares := core.middleWares
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	queryResult := handler(ctx, qc)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}
	return queryResult.Result.(*T), nil
}

func getMul[T any](ctx context.Context, core core, session session, qc *QueryContext) ([]*T, error) {
	var handler HandleFunc = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getMulHandler[T](ctx, core, session, qc)
	}
	middlewares := core.middleWares
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	queryResult := handler(ctx, qc)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}
	return queryResult.Result.([]*T), nil
}

func getHandler[T any](ctx context.Context, session session, core core, qc *QueryContext) *QueryResult {
	query, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	rows, err := session.queryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	if !rows.Next() {
		return &QueryResult{
			Err: errors.New("not data"),
		}
	}

	tp := new(T)
	tableModel, err := core.r.Get(tp)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	val := core.creator(tp, tableModel)
	err = val.SetColumns(rows)
	return &QueryResult{
		Result: tp,
		Err:    err,
	}
}

func getMulHandler[T any](ctx context.Context, core core, session session, qc *QueryContext) *QueryResult {
	query, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	rows, err := session.queryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	tpArr := make([]*T, 0)
	for rows.Next() {
		tp := new(T)
		tableModel, err := core.r.Get(tp)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		val := core.creator(tp, tableModel)
		err = val.SetColumns(rows)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		tpArr = append(tpArr, tp)
	}
	return &QueryResult{
		Result: tpArr,
	}
}
