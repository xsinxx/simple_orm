package accessLog

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/simple_orm"
)

type LogMiddleWare struct {
}

func NewLogMiddleWare() *LogMiddleWare {
	return &LogMiddleWare{}
}

func (l *LogMiddleWare) printLog(query *simple_orm.Query) {
	spew.Println(query)
}

func (l *LogMiddleWare) Build() simple_orm.MiddleWare {
	return func(next simple_orm.HandleFunc) simple_orm.HandleFunc {
		return func(ctx context.Context, qc *simple_orm.QueryContext) *simple_orm.QueryResult {
			query, err := qc.Builder.Build()
			if err != nil {
				return &simple_orm.QueryResult{
					Err: err,
				}
			}
			l.printLog(query)
			return next(ctx, qc)
		}
	}
}
