package simple_orm

import "github.com/simple_orm/model"

// Column 列名
type Column struct {
	name string
}

func (c *Column) expr() {

}

func NewColumn(name string) *Column {
	return &Column{
		name: name,
	}
}

// 列的右侧可能是表达式
func exprOf(e any) model.Expression {
	switch exp := e.(type) {
	case model.Expression:
		return exp
	default:
		return NewValue(exp)
	}
}

func (c *Column) LT(val any) *Predicate {
	return &Predicate{
		left:  c,
		right: exprOf(val),
		op:    model.opLT,
	}
}

func (c *Column) EQ(val any) *Predicate {
	return &Predicate{
		left:  c,
		right: exprOf(val),
		op:    model.opEQ,
	}
}

func (c *Column) GT(val any) *Predicate {
	return &Predicate{
		left:  c,
		right: exprOf(val),
		op:    model.opGT,
	}
}
