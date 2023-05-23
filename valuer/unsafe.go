package valuer

import (
	"database/sql"
	"errors"
	"github.com/simple_orm"
	"reflect"
	"unsafe"
)

type UnsafeValue struct {
	address    unsafe.Pointer
	tableModel *simple_orm.TableModel
}

func NewUnsafeValue(address any, meta *simple_orm.TableModel) Value {
	return UnsafeValue{
		address:    unsafe.Pointer(reflect.ValueOf(address).Pointer()),
		tableModel: meta,
	}
}

func (u UnsafeValue) SetColumns(rows *sql.Rows) error {
	columnsFromDB, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(columnsFromDB) != len(u.tableModel.Col2Field) {
		return errors.New("the number of columns doesn't match")
	}
	colVal := make([]interface{}, 0, len(columnsFromDB))
	for i, column := range columnsFromDB {
		field, ok := u.tableModel.Col2Field[column]
		if !ok {
			return errors.New("column from db not contained in meta")
		}
		ptr := unsafe.Pointer(uintptr(u.address) + field.Offset)
		val := reflect.NewAt(field.Typ, ptr) // 将指针写入对应的地址写入值
		colVal[i] = val.Interface()
	}
	return rows.Scan(colVal)
}
