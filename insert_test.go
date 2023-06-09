package simple_orm

import (
	"errors"
	"github.com/simple_orm/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	db := memoryDB4UnitTest(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 一个都不插入
			name:    "no value",
			q:       NewInserter[model.TestModel](db).Values(),
			wantErr: errors.New("[insert] insert zero row"),
		},
		{
			name: "single values",
			q: NewInserter[model.TestModel](db).Values(
				&model.TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
				}),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`id`,`first_name`,`age`) VALUES(?,?,?);",
				Args: []any{int64(1), "Deng", int8(18)},
			},
		},
		{
			name: "multiple values",
			q: NewInserter[model.TestModel](db).Values(
				&model.TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
				},
				&model.TestModel{
					Id:        2,
					FirstName: "Da",
					Age:       19,
				}),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`id`,`first_name`,`age`) VALUES(?,?,?),(?,?,?);",
				Args: []any{int64(1), "Deng", int8(18), int64(2), "Da", int8(19)},
			},
		},
		{
			// 指定列
			name: "specify columns",
			q: NewInserter[model.TestModel](db).Values(
				&model.TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
				}).Columns("FirstName", "Age"),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`first_name`,`age`) VALUES(?,?);",
				Args: []any{"Deng", int8(18)},
			},
		},
		{
			// upsert
			name: "upsert",
			q: NewInserter[model.TestModel](db).Values(
				&model.TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
				}).OnDuplicateKey().Update(Assign("FirstName", "Da")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`) VALUES(?,?,?) " +
					"ON DUPLICATE KEY UPDATE `first_name`=?;",
				Args: []any{int64(1), "Deng", int8(18), "Da"},
			},
		},
		{
			// upsert invalid column
			name: "upsert invalid column",
			q: NewInserter[model.TestModel](db).Values(
				&model.TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
				}).OnDuplicateKey().Update(Assign("Invalid", "Da")),
			wantErr: errors.New("column name not exists"),
		},
		{
			// 使用原本插入的值
			name: "upsert use insert value",
			q: NewInserter[model.TestModel](db).Values(
				&model.TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
				},
				&model.TestModel{
					Id:        2,
					FirstName: "Da",
					Age:       19,
				}).OnDuplicateKey().Update(NewColumn("FirstName")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`) VALUES(?,?,?),(?,?,?) " +
					"ON DUPLICATE KEY UPDATE `first_name`=VALUES(`first_name`);",
				Args: []any{int64(1), "Deng", int8(18), int64(2), "Da", int8(19)},
			},
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
