package valuer

import (
	"database/sql"
	"github.com/simple_orm/model"
)

// Value 是对结构体实例的内部抽象
type Value interface {
	// GetValByColName 获取colName在结构体中对应字段的值，GetValByColName(age)返回的结果是具体的年龄，如25
	GetValByColName(colName string) (any, error)
	// SetColumns 设置新值
	SetColumns(rows *sql.Rows) error
}

type Creator func(val interface{}, meta *model.TableModel) Value
