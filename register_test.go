package simple_orm

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDB_Option(t *testing.T) {
	models := map[reflect.Type]*TableModel{
		reflect.TypeOf("user"): &TableModel{
			TableName: "user_t",
		},
	}
	r := &Registry{
		models: models,
	}
	db := MustNewDB(DBWithRegister(r))
	assert.Equal(t, db.r, r)
}

func TestRegistry_get(t *testing.T) {
	testCases := []struct {
		name      string
		val       any
		wantModel *TableModel
		wantErr   error
	}{
		{
			// 指针
			name: "pointer",
			val:  &TestModel{},
			wantModel: &TableModel{
				TableName: "test_model",
				Col2Field: map[string]*field{
					"Id": {
						ColumnName: "id",
					},
					"FirstName": {
						ColumnName: "first_name",
					},
					"Age": {
						ColumnName: "age",
					},
					"LastName": {
						ColumnName: "last_name",
					},
				},
			},
		},
	}

	r := &Registry{
		models: map[reflect.Type]*TableModel{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}
