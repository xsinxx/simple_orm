package simple_orm

import "github.com/simple_orm/model"

// Expression 顶层抽象，接口的实现类有表达式、列、值。
// eg：(`Age` > 13) AND (`Age` < 24)，(`Age` > 13)是表达式、Age是列、13填写的是值
type Expression interface {
	expr()
}

// Predicate 表达式
type Predicate struct {
	left  Expression
	right Expression
	op    model.Op
}

func (p *Predicate) expr() {

}

func (p *Predicate) And(right *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    model.OpAnd,
		right: right,
	}
}

func (p *Predicate) Or(right *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    model.OpOr,
		right: right,
	}
}

func Not(right *Predicate) *Predicate {
	return &Predicate{
		op:    model.OpNot,
		right: right,
	}
}
