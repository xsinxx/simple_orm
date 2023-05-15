package simple_orm

// Predicate 表达式
type Predicate struct {
	left  Expression
	right Expression
	op    op
}

func (p *Predicate) expr() {

}

func (p *Predicate) And(right *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    opAnd,
		right: right,
	}
}

func (p *Predicate) Or(right *Predicate) *Predicate {
	return &Predicate{
		left:  p,
		op:    opOr,
		right: right,
	}
}

func Not(right *Predicate) *Predicate {
	return &Predicate{
		op:    opNot,
		right: right,
	}
}
