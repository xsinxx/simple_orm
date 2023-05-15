package simple_orm

import "database/sql"

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
