package simple_orm

import "github.com/simple_orm/model"

// Predicate 表达式
type Predicate struct {
	left  model.Expression
	right model.Expression
	op    model.op
}

func (p *Predicate) expr() {

}

func (p *Predicate) And(right *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    model.opAnd,
		right: right,
	}
}

func (p *Predicate) Or(right *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    model.opOr,
		right: right,
	}
}

func Not(right *Predicate) *Predicate {
	return &Predicate{
		op:    model.opNot,
		right: right,
	}
}
