package valuer

import (
	"database/sql"
	"github.com/simple_orm"
)

// Value 是对结构体实例的内部抽象
type Value interface {
	// SetColumns 设置新值
	SetColumns(rows *sql.Rows) error
}

type Creator func(val interface{}, meta *simple_orm.TableModel) Value
