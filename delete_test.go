package simple_orm

import (
	"errors"
	"github.com/simple_orm/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleter_Build(t *testing.T) {
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// From 都不调用
			name: "no from",
			q:    NewDeleter[model.TestModel](memoryDB4UnitTest(t)),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			// 调用 FROM
			name: "with from",
			q:    NewDeleter[model.TestModel](memoryDB4UnitTest(t)).From("`test_model_t`"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model_t`;",
			},
		},
		{
			// 调用 FROM，但是传入空字符串
			name: "empty from",
			q:    NewDeleter[model.TestModel](memoryDB4UnitTest(t)).From(""),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			// 调用 FROM，同时出入看了 DB
			name: "with db",
			q:    NewDeleter[model.TestModel](memoryDB4UnitTest(t)).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_db`.`test_model`;",
			},
		},
		{
			// 单一简单条件
			name: "single and simple predicate",
			q: NewDeleter[model.TestModel](memoryDB4UnitTest(t)).From("`test_model_t`").
				Where(NewColumn("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model_t` WHERE `id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			q: NewDeleter[model.TestModel](memoryDB4UnitTest(t)).
				Where(NewColumn("Age").GT(18), NewColumn("Age").LT(35)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 AND
			name: "and",
			q: NewDeleter[model.TestModel](memoryDB4UnitTest(t)).
				Where(NewColumn("Age").GT(18).And(NewColumn("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name: "or",
			q: NewDeleter[model.TestModel](memoryDB4UnitTest(t)).
				Where(NewColumn("Age").GT(18).Or(NewColumn("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`age` > ?) OR (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name: "not",
			q:    NewDeleter[model.TestModel](memoryDB4UnitTest(t)).Where(Not(NewColumn("Age").GT(18))),
			wantQuery: &Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "DELETE FROM `test_model` WHERE  NOT (`age` > ?);",
				Args: []any{18},
			},
		},
		{
			// 非法列
			name:    "invalid column",
			q:       NewDeleter[model.TestModel](memoryDB4UnitTest(t)).Where(Not(NewColumn("Invalid").GT(18))),
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
