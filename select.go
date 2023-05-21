package simple_orm

import (
	"errors"
	"strings"
)

// Selector 用于构造 SELECT 语句
type Selector[T any] struct {
	sb          strings.Builder
	table       string
	where       []*Predicate
	args        []any
	tableModels *tableModel
	db          *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}

// From 指定表名，如果是空字符串，那么将会使用默认表名
func (s *Selector[T]) From(tbl string) *Selector[T] {
	s.table = tbl
	return s
}

func (s *Selector[T]) Where(where ...*Predicate) *Selector[T] {
	s.where = where
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	s.tableModels, err = s.db.r.get(t)
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT * FROM ")
	// 表名
	if s.table == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.tableModels.tableName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.table)
	}
	// where
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		err := s.buildExpression(p)
		if err != nil {
			return nil, err
		}
	}

	s.sb.WriteString(";")
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

// 递归解析表达式
// (`Age` > 13) AND (`Age` < 24)
func (s *Selector[T]) buildExpression(e Expression) error {
	switch expr := e.(type) {
	case *Column: // 列， eg：`Age`
		if _, ok := s.tableModels.fieldMap[expr.name]; !ok {
			return errors.New("illegal field")
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(s.tableModels.fieldMap[expr.name].columnName)
		s.sb.WriteByte('`')
	case *Value: // 值，eg： 13
		s.sb.WriteByte('?')
		s.args = append(s.args, expr.val)
	case *Predicate: // 表达式
		// 左侧表达式
		_, lp := expr.left.(*Predicate)
		if lp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(expr.left); err != nil {
			return err
		}
		if lp {
			s.sb.WriteByte(')')
		}
		// 链接符
		s.sb.WriteByte(' ')
		s.sb.WriteString(string(expr.op))
		s.sb.WriteByte(' ')
		// 右侧表达式
		_, rp := expr.right.(*Predicate)
		if rp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(expr.right); err != nil {
			return err
		}
		if rp {
			s.sb.WriteByte(')')
		}
	}
	return nil
}
