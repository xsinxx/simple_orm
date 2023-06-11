package access_log

import (
	"context"
	"fmt"
	"github.com/simple_orm"
)

type LogMiddleWare struct {
}

func NewLogMiddleWare() *LogMiddleWare {
	return &LogMiddleWare{}
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
			fmt.Println(query.SQL)
			return next(ctx, qc)
		}
	}
}
