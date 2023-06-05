package simple_orm

import (
	"errors"
)

var (
	mySQLDialect = &MySQLDialect{}
)

type Dialect interface {
	// Upsert 不同的数据库实现不同的Upsert
	Upsert(builder *Builder, upsert *UpsertKey) error
}

type standardSQL struct {
}

type MySQLDialect struct {
	standardSQL
}

func (m *MySQLDialect) Upsert(builder *Builder, upsert *UpsertKey) error {
	builder.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for _, assign := range upsert.assigns {
		switch e := assign.(type) {
		case *Column:
			columnName := e.name
			field, ok := builder.tableModels.Col2Field[columnName]
			if !ok {
				return errors.New("column name not exists")
			}
			// 详细结构参照单测
			builder.sb.WriteString("`")
			builder.sb.WriteString(field.ColumnName)
			builder.sb.WriteString("`=VALUES(`")
			builder.sb.WriteString(field.ColumnName)
			builder.sb.WriteString("`)")
		case *Assignment:
			columnName := e.ColumnName
			field, ok := builder.tableModels.Col2Field[columnName]
			if !ok {
				return errors.New("column name not exists")
			}
			builder.sb.WriteString("`")
			builder.sb.WriteString(field.ColumnName)
			builder.sb.WriteString("`")
			builder.sb.WriteString("=?")
			builder.args = append(builder.args, e.Val)
		}
	}
	return nil
}
