package model

import (
	"database/sql"
	"reflect"
	"sync"
)

type Query struct {
	SQL  string
	Args []any
}

type QueryBuilder interface {
	Build() (*Query, error)
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

// Expression 顶层抽象，接口的实现类有表达式、列、值。
// eg：(`Age` > 13) AND (`Age` < 24)，(`Age` > 13)是表达式、Age是列、13填写的是值
type Expression interface {
	expr()
}

type op string

const (
	opAnd = "AND"
	opOr  = "OR"
	opNot = "NOT"
	opLT  = "<"
	opGT  = ">"
	opEQ  = "="
)

type Field struct {
	ColumnName string // 对应的数据库中表的列
	Typ        reflect.Type
	Offset     uintptr
}

type TableModel struct {
	TableName string            // 表名
	Tag2Field map[string]*Field // 标签名到字段的映射
	Col2Field map[string]*Field // 列名到字段的映射
}

// Registry 注册中心，存储表信息
type Registry struct {
	lock        sync.RWMutex // 防止读写冲突
	TableModels map[reflect.Type]*TableModel
}
