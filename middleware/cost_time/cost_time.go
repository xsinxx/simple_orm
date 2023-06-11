package cost_time

import (
	"context"
	"fmt"
	"github.com/simple_orm"
	"time"
)

type CostTimeMiddleWare struct {
}

func NewCostTimeMiddleWare() *CostTimeMiddleWare {
	return &CostTimeMiddleWare{}
}

func (s *CostTimeMiddleWare) Build() simple_orm.MiddleWare {
	return func(next simple_orm.HandleFunc) simple_orm.HandleFunc {
		return func(ctx context.Context, qc *simple_orm.QueryContext) *simple_orm.QueryResult {
			defer func(now time.Time) {
				fmt.Printf("query cost time:%d\n", time.Since(now).Milliseconds())
			}(time.Now())
			return next(ctx, qc)
		}
	}
}
