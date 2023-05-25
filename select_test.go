package simple_orm

import (
	"errors"
	"github.com/simple_orm/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	testCases := []struct {
		name      string
		q         model.QueryBuilder
		wantQuery *model.Query
		wantErr   error
	}{
		{
			// From 都不调用
			name: "no from",
			q:    NewSelector[model.TestModel](MustNewDB()),
			wantQuery: &model.Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM
			name: "with from",
			q:    NewSelector[model.TestModel](MustNewDB()).From("`test_model_t`"),
			wantQuery: &model.Query{
				SQL: "SELECT * FROM `test_model_t`;",
			},
		},
		{
			// 调用 FROM，但是传入空字符串
			name: "empty from",
			q:    NewSelector[model.TestModel](MustNewDB()).From(""),
			wantQuery: &model.Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM，同时出入看了 DB
			name: "with db",
			q:    NewSelector[model.TestModel](MustNewDB()).From("`test_db`.`test_model`"),
			wantQuery: &model.Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		},
		{
			// 单一简单条件
			name: "single and simple predicate",
			q: NewSelector[model.TestModel](MustNewDB()).From("`test_model_t`").
				Where(NewColumn("Id").EQ(1)),
			wantQuery: &model.Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			q: NewSelector[model.TestModel](MustNewDB()).
				Where(NewColumn("Age").GT(18), NewColumn("Age").LT(35)),
			wantQuery: &model.Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 AND
			name: "and",
			q: NewSelector[model.TestModel](MustNewDB()).
				Where(NewColumn("Age").GT(18).And(NewColumn("Age").LT(35))),
			wantQuery: &model.Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name: "or",
			q: NewSelector[model.TestModel](MustNewDB()).
				Where(NewColumn("Age").GT(18).Or(NewColumn("Age").LT(35))),
			wantQuery: &model.Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) OR (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name: "not",
			q:    NewSelector[model.TestModel](MustNewDB()).Where(Not(NewColumn("Age").GT(18))),
			wantQuery: &model.Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` > ?);",
				Args: []any{18},
			},
		},
		{
			// 非法列
			name:    "invalid column",
			q:       NewSelector[model.TestModel](MustNewDB()).Where(Not(NewColumn("Invalid").GT(18))),
			wantErr: errors.New("illegal field"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
