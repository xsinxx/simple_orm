package simple_orm

import (
	"context"
	"errors"
)

// Selector 用于构造 SELECT 语句
type Selector[T any] struct {
	Builder         // builder是Insert & Select公共部分
	core            // core中是元数据信息
	session session // session是db或tx
	table   string
	where   []*Predicate
	groupBy []*Column
	having  *Predicate
	orderBy []*OrderBy
	limit   int
	offset  int
}

func NewSelector[T any](session session) *Selector[T] {
	return &Selector[T]{
		core:    session.getCore(),
		session: session,
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

func (s *Selector[T]) GroupBy(columns ...*Column) *Selector[T] {
	s.groupBy = columns
	return s
}

func (s *Selector[T]) Having(having *Predicate) *Selector[T] {
	s.having = having
	return s
}

func (s *Selector[T]) OrderBy(columns ...*OrderBy) *Selector[T] {
	s.orderBy = columns
	return s
}

func (s *Selector[T]) Limit(limit int) *Selector[T] {
	s.limit = limit
	return s
}

func (s *Selector[T]) Offset(offset int) *Selector[T] {
	s.offset = offset
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	s.tableModels, err = s.r.Get(t)
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT * FROM ")
	// table aggregateFunction
	if s.table == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.tableModels.TableName)
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
		err = s.buildExpression(p)
		if err != nil {
			return nil, err
		}
	}

	// group by
	if len(s.groupBy) > 0 {
		s.sb.WriteString(" GROUP BY ")
		for i, v := range s.groupBy {
			err = s.buildExpression(v)
			if err != nil {
				return nil, err
			}
			if i != len(s.groupBy)-1 {
				s.sb.WriteString(",")
			}
		}
	}

	// having
	if s.having != nil {
		if len(s.groupBy) == 0 {
			return nil, errors.New("[having] group by clause is not exists")
		}
		s.sb.WriteString(" HAVING ")
		err = s.buildExpression(s.having)
		if err != nil {
			return nil, err
		}
	}

	// order by
	if len(s.orderBy) > 0 {
		s.sb.WriteString(" ORDER BY ")
		for i, v := range s.orderBy {
			field, ok := s.tableModels.Col2Field[v.name]
			if !ok {
				return nil, errors.New("illegal field")
			}
			s.sb.WriteString("`" + field.ColumnName + "` ")
			s.sb.WriteString(string(v.order))
			if i != len(s.orderBy)-1 {
				s.sb.WriteString(",")
			}
		}
	}

	// limit
	if s.limit != 0 {
		s.sb.WriteString(" LIMIT ?")
		s.args = append(s.args, s.limit)
	}

	// offset
	if s.offset != 0 {
		s.sb.WriteString(" OFFSET ?")
		s.args = append(s.args, s.offset)
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
	case *Aggregate:
		field, ok := s.tableModels.Col2Field[expr.name]
		if !ok {
			return errors.New("illegal field")
		}
		s.sb.WriteString(string(expr.aggregateFunction))
		s.sb.WriteString("(`")
		s.sb.WriteString(field.ColumnName)
		s.sb.WriteString("`)")
	case *Column: // 列， eg：`Age`
		if _, ok := s.tableModels.Col2Field[expr.name]; !ok {
			return errors.New("illegal field")
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(s.tableModels.Col2Field[expr.name].ColumnName)
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

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	var handler HandleFunc = s.getHandler
	middlewares := s.middleWares
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	qc := &QueryContext{}
	queryResult := handler(ctx, qc)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}
	return queryResult.Result.(*T), nil
}

func (s *Selector[T]) GetMul(ctx context.Context) ([]*T, error) {
	var handler HandleFunc = s.getMulHandler
	middlewares := s.middleWares
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	qc := &QueryContext{}
	queryResult := handler(ctx, qc)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}
	return queryResult.Result.([]*T), nil
}

func (s *Selector[T]) getHandler(ctx context.Context, qc *QueryContext) *QueryResult {
	query, err := s.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	rows, err := s.session.queryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	if !rows.Next() {
		return &QueryResult{
			Err: errors.New("not data"),
		}
	}

	tp := new(T)
	tableModel, err := s.r.Get(tp)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	val := s.creator(tp, tableModel)
	err = val.SetColumns(rows)
	return &QueryResult{
		Result: tp,
		Err:    err,
	}
}

func (s *Selector[T]) getMulHandler(ctx context.Context, qc *QueryContext) *QueryResult {
	query, err := s.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	rows, err := s.session.queryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	tpArr := make([]*T, 0)
	for rows.Next() {
		tp := new(T)
		tableModel, err := s.r.Get(tp)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		val := s.creator(tp, tableModel)
		err = val.SetColumns(rows)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		tpArr = append(tpArr, tp)
	}
	return &QueryResult{
		Result: tpArr,
	}
}
