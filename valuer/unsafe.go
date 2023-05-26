package valuer

import (
	"database/sql"
	"github.com/simple_orm/model"
	"reflect"
	"unsafe"
)

type UnsafeValue struct {
	address    unsafe.Pointer
	tableModel *model.TableModel
}

func NewUnsafeValue(address any, meta *model.TableModel) Value {
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
	if len(columnsFromDB) > len(u.tableModel.Col2Field) {
		return ColumnsNotMatch
	}
	colVal := make([]interface{}, len(columnsFromDB))
	for i, column := range columnsFromDB {
		field, ok := u.tableModel.Col2Field[column]
		if !ok {
			return ColumnsNotExists
		}
		ptr := unsafe.Pointer(uintptr(u.address) + field.Offset)
		val := reflect.NewAt(field.Typ, ptr) // 将指针对应的地址写入值
		colVal[i] = val.Interface()
	}
	return rows.Scan(colVal...)
}
