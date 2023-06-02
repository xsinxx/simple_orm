package simple_orm

type Assignable interface {
	assign()
}

type Assignment struct {
	ColumnName string
	Val        any
}

func Assign(column string, val any) *Assignment {
	return &Assignment{
		ColumnName: column,
		Val:        val,
	}
}

func (a *Assignment) assign() {}
