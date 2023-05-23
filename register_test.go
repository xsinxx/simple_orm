package simple_orm

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDB_Option(t *testing.T) {
	models := map[reflect.Type]*tableModel{
		reflect.TypeOf("user"): &tableModel{
			tableName: "user_t",
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
		wantModel *tableModel
		wantErr   error
	}{
		{
			// 指针
			name: "pointer",
			val:  &TestModel{},
			wantModel: &tableModel{
				tableName: "test_model",
				fieldMap: map[string]*field{
					"Id": {
						columnName: "id",
					},
					"FirstName": {
						columnName: "first_name",
					},
					"Age": {
						columnName: "age",
					},
					"LastName": {
						columnName: "last_name",
					},
				},
			},
		},
	}

	r := &Registry{
		models: map[reflect.Type]*tableModel{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}
