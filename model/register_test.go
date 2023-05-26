package model

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRegistry_Get(t *testing.T) {
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
				Col2Field: map[string]*Field{
					"Id": {
						ColumnName: "id",
						Typ:        reflect.TypeOf(int64(0)),
						TypName:    "Id",
						Offset:     0,
					},
					"FirstName": {
						ColumnName: "first_name",
						Typ:        reflect.TypeOf(""),
						TypName:    "FirstName",
						Offset:     8,
					},
					"Age": {
						ColumnName: "age",
						Typ:        reflect.TypeOf(int8(0)),
						TypName:    "Age",
						Offset:     24,
					},
				},
				Tag2Field: map[string]*Field{
					"Id": {
						ColumnName: "id",
						Typ:        reflect.TypeOf(int64(0)),
						TypName:    "Id",
						Offset:     0,
					},
					"FirstName": {
						ColumnName: "first_name",
						Typ:        reflect.TypeOf(""),
						TypName:    "FirstName",
						Offset:     8,
					},
					"Age": {
						ColumnName: "age",
						Typ:        reflect.TypeOf(int8(0)),
						TypName:    "Age",
						Offset:     24,
					},
				},
			},
		},
		{
			// 多级指针
			name: "pointer",
			val: func() any {
				t := &TestModel{}
				return &t
			}(),
			wantModel: &TableModel{
				TableName: "test_model",
				Col2Field: map[string]*Field{
					"Id": {
						ColumnName: "id",
						Typ:        reflect.TypeOf(int64(0)),
						TypName:    "Id",
						Offset:     0,
					},
					"FirstName": {
						ColumnName: "first_name",
						Typ:        reflect.TypeOf(""),
						TypName:    "FirstName",
						Offset:     8,
					},
					"Age": {
						ColumnName: "age",
						Typ:        reflect.TypeOf(int8(0)),
						TypName:    "Age",
						Offset:     24,
					},
				},
				Tag2Field: map[string]*Field{
					"Id": {
						ColumnName: "id",
						Typ:        reflect.TypeOf(int64(0)),
						TypName:    "Id",
						Offset:     0,
					},
					"FirstName": {
						ColumnName: "first_name",
						Typ:        reflect.TypeOf(""),
						TypName:    "FirstName",
						Offset:     8,
					},
					"Age": {
						ColumnName: "age",
						Typ:        reflect.TypeOf(int8(0)),
						TypName:    "Age",
						Offset:     24,
					},
				},
			},
		},
		{
			name:    "map",
			val:     map[string]string{},
			wantErr: errors.New("type is wrong"),
		},
		{
			name: "add Tag",
			val: func() any {
				type Tag struct {
					Level int64 `orm:"identity"`
				}
				return &Tag{}
			}(),
			wantModel: &TableModel{
				TableName: "tag",
				Tag2Field: map[string]*Field{
					"identity": {
						ColumnName: "level",
						Typ:        reflect.TypeOf(int64(0)),
						TypName:    "Level",
						Offset:     0,
					},
				},
				Col2Field: map[string]*Field{
					"Level": {
						ColumnName: "level",
						Typ:        reflect.TypeOf(int64(0)),
						TypName:    "Level",
						Offset:     0,
					},
				},
			},
		},
	}

	r := &Registry{
		TableModels: map[reflect.Type]*TableModel{},
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
