package simple_orm

// Value 值
type Value struct {
	val any
}

func (v *Value) expr() {
}

func NewValue(val any) *Value {
	return &Value{
		val: val,
	}
}
