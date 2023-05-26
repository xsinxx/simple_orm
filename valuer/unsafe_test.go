package valuer

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simple_orm/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_unsafeValue_SetColumn(t *testing.T) {
	testCases := []struct {
		name    string
		cs      map[string][]byte
		val     *model.TestModel
		wantVal *model.TestModel
		wantErr error
	}{
		{
			name: "invalid field",
			cs: map[string][]byte{
				"invalid_column": nil,
			},
			wantErr: ColumnsNotExists,
		},
		{
			name: "normal deal result set",
			cs: map[string][]byte{
				"Id":        []byte("9426"),
				"FirstName": []byte("zhu zhu"),
				"Age":       []byte("66"),
			},
			val: &model.TestModel{},
			wantVal: &model.TestModel{
				Id:        9426,
				FirstName: "zhu zhu",
				Age:       66,
			},
		},
		{
			name: "normal deal result set",
			cs: map[string][]byte{
				"Id":        []byte("9426"),
				"FirstName": []byte("zhu zhu"),
				"Location":  []byte("BJ"),
			},
			val:     &model.TestModel{},
			wantErr: ColumnsNotExists,
		},
	}

	r := model.NewRegistry()
	meta, err := r.Get(&model.TestModel{})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			if err != nil {
				t.Fatal(err)
			}
			val := NewUnsafeValue(tc.val, meta)
			cols := make([]string, 0, len(tc.cs))
			colVals := make([]driver.Value, 0, len(tc.cs))
			for k, v := range tc.cs {
				cols = append(cols, k)
				colVals = append(colVals, v)
			}
			// 当db.Query执行ExpectQuery中的语句时，返回结果是WillReturnRows中写入的结果
			mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).AddRow(colVals...))
			rows, _ := db.Query("SELECT *")
			rows.Next()
			err = val.SetColumns(rows)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			if tc.wantErr != nil {
				t.Fatalf("期望得到错误，但是并没有得到 %v", tc.wantErr)
			}
			assert.Equal(t, tc.wantVal, tc.val)
		})
	}
}
