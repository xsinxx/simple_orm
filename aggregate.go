package simple_orm

import "github.com/simple_orm/model"

// Aggregate 聚合函数，eg: SUM、AVG
type Aggregate struct {
	aggregateFunction AggregateFunction
	name              string
}

func (a *Aggregate) expr() {

}

func NewAggregate(name string, aggregateFunction AggregateFunction) *Aggregate {
	return &Aggregate{
		aggregateFunction: aggregateFunction,
		name:              name,
	}
}

func (a *Aggregate) LT(val any) *Predicate {
	return &Predicate{
		left:  a,
		right: exprOf(val),
		op:    model.OpLT,
	}
}

func (a *Aggregate) EQ(val any) *Predicate {
	return &Predicate{
		left:  a,
		right: exprOf(val),
		op:    model.OpEQ,
	}
}

func (a *Aggregate) GT(val any) *Predicate {
	return &Predicate{
		left:  a,
		right: exprOf(val),
		op:    model.OpGT,
	}
}
