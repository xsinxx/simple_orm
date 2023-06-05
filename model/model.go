package model

import (
	"reflect"
	"sync"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	//LastName  *sql.NullString
}

type Op string

const (
	OpAnd = "AND"
	OpOr  = "OR"
	OpNot = "NOT"
	OpLT  = "<"
	OpGT  = ">"
	OpEQ  = "="
)

type Field struct {
	ColumnName string // 对应的数据库中表的列
	Typ        reflect.Type
	TypName    string
	Offset     uintptr
}

type TableModel struct {
	TableName   string            // 表名
	Tag2Field   map[string]*Field // 标签名到字段的映射
	Col2Field   map[string]*Field // 列名到字段的映射
	ColumnNames []string          // 列名数组，由于map的遍历是乱序，因此用数组保证顺序
}

// Registry 注册中心，存储表信息
type Registry struct {
	lock        sync.RWMutex // 防止读写冲突
	TableModels map[reflect.Type]*TableModel
}
