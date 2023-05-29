package simple_orm

import (
	"errors"
	"github.com/simple_orm/model"
	"strings"
)

type Delete[T any] struct {
	values      []any
	table       string
	where       []*Predicate
	db          *DB
	args        []any
	sb          strings.Builder
	tableModels *model.TableModel
}

func NewDeleter[T any](db *DB) *Delete[T] {
	return &Delete[T]{
		db: db,
	}
}

// From 指定表名，如果是空字符串，那么将会使用默认表名
func (d *Delete[T]) From(tbl string) *Delete[T] {
	d.table = tbl
	return d
}

func (d *Delete[T]) Where(where ...*Predicate) *Delete[T] {
	d.where = where
	return d
}

func (d *Delete[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	d.tableModels, err = d.db.r.Get(t)
	if err != nil {
		return nil, err
	}
	d.sb.WriteString("DELETE FROM ")
	// table name
	if d.table == "" {
		d.sb.WriteByte('`')
		d.sb.WriteString(d.tableModels.TableName)
		d.sb.WriteByte('`')
	} else {
		d.sb.WriteString(d.table)
	}

	// where
	if len(d.where) > 0 {
		d.sb.WriteString(" WHERE ")
		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}
		err = d.buildExpression(p)
		if err != nil {
			return nil, err
		}
	}

	d.sb.WriteString(";")
	return &Query{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}

// 递归解析表达式
// (`Age` > 13) AND (`Age` < 24)
func (d *Delete[T]) buildExpression(e Expression) error {
	switch expr := e.(type) {
	case *Aggregate:
		field, ok := d.tableModels.Col2Field[expr.name]
		if !ok {
			return errors.New("illegal field")
		}
		d.sb.WriteString(string(expr.aggregateFunction))
		d.sb.WriteString("(`")
		d.sb.WriteString(field.ColumnName)
		d.sb.WriteString("`)")
	case *Column: // 列， eg：`Age`
		if _, ok := d.tableModels.Col2Field[expr.name]; !ok {
			return errors.New("illegal field")
		}
		d.sb.WriteByte('`')
		d.sb.WriteString(d.tableModels.Col2Field[expr.name].ColumnName)
		d.sb.WriteByte('`')
	case *Value: // 值，eg： 13
		d.sb.WriteByte('?')
		d.args = append(d.args, expr.val)
	case *Predicate: // 表达式
		// 左侧表达式
		_, lp := expr.left.(*Predicate)
		if lp {
			d.sb.WriteByte('(')
		}
		if err := d.buildExpression(expr.left); err != nil {
			return err
		}
		if lp {
			d.sb.WriteByte(')')
		}
		// 链接符
		d.sb.WriteByte(' ')
		d.sb.WriteString(string(expr.op))
		d.sb.WriteByte(' ')
		// 右侧表达式
		_, rp := expr.right.(*Predicate)
		if rp {
			d.sb.WriteByte('(')
		}
		if err := d.buildExpression(expr.right); err != nil {
			return err
		}
		if rp {
			d.sb.WriteByte(')')
		}
	}
	return nil
}