package valuer

import (
	"database/sql"
	"github.com/simple_orm/model"
	"reflect"
)

type ReflectValue struct {
	val  reflect.Value
	meta *model.TableModel
}

func NewReflectValue(val any, meta *model.TableModel) Value {
	return ReflectValue{
		val:  reflect.ValueOf(val).Elem(),
		meta: meta,
	}
}

// SetColumns
//  r.val只能是结构体指针，目标是将该结构体中的字段set数据库中读取出的值
//  ==> set的数据只能是reflect.Value类型，因此需要从rows中读取出reflect.Value
//  ==> rows.Scan只能接收[]interface{}，这个[]interface{}是指针数组
//  ==> 由于是reflect.New返回结果是指针，在调用scan时会同时更新colValues & colEleValues
///*
func (r ReflectValue) SetColumns(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(cols) != len(r.meta.Col2Field) {
		return ColumnsNotMatch
	}
	colValues := make([]interface{}, len(cols))
	colEleValues := make([]reflect.Value, len(cols))
	for i, col := range cols {
		field, ok := r.meta.Col2Field[col]
		if !ok {
			return ColumnsNotExists
		}
		ptr := reflect.New(field.Typ)  // new一个空指针
		colValues[i] = ptr.Interface() // 指针中存储的是地址，通过Interface()取的是地址
		colEleValues[i] = ptr.Elem()
	}
	// scan中传入的值是指针中存储的地址
	if err = rows.Scan(colValues); err != nil {
		return err
	}
	for i, col := range cols {
		field := r.meta.Col2Field[col]
		fd := r.val.FieldByName(field.ColumnName)
		fd.Set(colEleValues[i])
	}
	return nil
}
