package simple_orm

// OrderBy 聚合函数，eg: SUM、AVG
type OrderBy struct {
	name  string
	order Order
}

func (o *OrderBy) expr() {

}

func Asc(name string) *OrderBy {
	return &OrderBy{
		order: ASCOrder,
		name:  name,
	}
}

func Desc(name string) *OrderBy {
	return &OrderBy{
		order: DESCOrder,
		name:  name,
	}
}
