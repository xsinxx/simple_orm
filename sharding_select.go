package simple_orm

import (
	"errors"
	"github.com/simple_orm/sharding"
	"github.com/valyala/bytebufferpool"
)

type ShardingSelector[T any] struct {
	ShardingBuilder         // shardingBuilder是SQL的公共部分
	core                    // core中是元数据信息
	session         session // session是db或tx
	table           string
	where           []*Predicate
	groupBy         []*Column
	having          *Predicate
	orderBy         []*OrderBy
	limit           int
	offset          int
}

func NewShardingSelector[T any](session session) *ShardingSelector[T] {
	return &ShardingSelector[T]{
		core:    session.getCore(),
		session: session,
	}
}

func (s *ShardingSelector[T]) Build() ([]*Query, error) {
	var (
		t   T
		err error
	)
	// tableModel
	tableModel, err := s.r.Get(t)
	if err != nil {
		return nil, err
	}
	s.tableModels = tableModel
	if s.core.algorithm == nil {
		return []*Query{}, errors.New("no valid algorithm")
	}
	// algorithm
	s.ShardingBuilder.algorithm = s.core.algorithm
	// byte buffer pool：申请一块内存
	s.stringBuffer = bytebufferpool.Get()
	dataSources, err := s.FindDataSource(s.where...)
	queries := make([]*Query, 0)
	// 构建单独的query & 将buffer归还
	// refer https://segmentfault.com/a/1190000039969499
	defer bytebufferpool.Put(s.stringBuffer)
	for _, dataSource := range dataSources {
		query, err := s.buildQuery(dataSource)
		if err != nil {
			return []*Query{}, err
		}
		queries = append(queries, query)
		s.stringBuffer.Reset()
	}
	return queries, nil
}

func (s *ShardingSelector[T]) buildQuery(dataSource *sharding.DataSource) (*Query, error) {
	var (
		err error
	)
	s.stringBuffer.WriteString("SELECT * FROM ")
	s.stringBuffer.WriteString(dataSource.DB + "." + dataSource.Table)
	// where
	if len(s.where) > 0 {
		s.stringBuffer.WriteString(" WHERE ")
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
		s.stringBuffer.WriteString(" GROUP BY ")
		for i, v := range s.groupBy {
			err = s.buildExpression(v)
			if err != nil {
				return nil, err
			}
			if i != len(s.groupBy)-1 {
				s.stringBuffer.WriteString(",")
			}
		}
	}

	// having
	if s.having != nil {
		if len(s.groupBy) == 0 {
			return nil, errors.New("[having] group by clause is not exists")
		}
		s.stringBuffer.WriteString(" HAVING ")
		err = s.buildExpression(s.having)
		if err != nil {
			return nil, err
		}
	}

	// order by
	if len(s.orderBy) > 0 {
		s.stringBuffer.WriteString(" ORDER BY ")
		for i, v := range s.orderBy {
			field, ok := s.tableModels.Col2Field[v.name]
			if !ok {
				return nil, errors.New("illegal field")
			}
			s.stringBuffer.WriteString("`" + field.ColumnName + "` ")
			s.stringBuffer.WriteString(string(v.order))
			if i != len(s.orderBy)-1 {
				s.stringBuffer.WriteString(",")
			}
		}
	}

	// limit
	if s.limit != 0 {
		s.stringBuffer.WriteString(" LIMIT ?")
		s.args = append(s.args, s.limit)
	}

	// offset
	if s.offset != 0 {
		s.stringBuffer.WriteString(" OFFSET ?")
		s.args = append(s.args, s.offset)
	}
	s.stringBuffer.WriteString(";")
	return &Query{
		SQL:  s.stringBuffer.String(),
		Args: s.args,
	}, nil
}

func (s *ShardingSelector[T]) buildExpression(e Expression) error {
	switch expr := e.(type) {
	case *Aggregate:
		field, ok := s.tableModels.Col2Field[expr.name]
		if !ok {
			return errors.New("illegal field")
		}
		s.stringBuffer.WriteString(string(expr.aggregateFunction))
		s.stringBuffer.WriteString("(`")
		s.stringBuffer.WriteString(field.ColumnName)
		s.stringBuffer.WriteString("`)")
	case *Column: // 列， eg：`Age`
		if _, ok := s.tableModels.Col2Field[expr.name]; !ok {
			return errors.New("illegal field")
		}
		s.stringBuffer.WriteByte('`')
		s.stringBuffer.WriteString(s.tableModels.Col2Field[expr.name].ColumnName)
		s.stringBuffer.WriteByte('`')
	case *Value: // 值，eg： 13
		s.stringBuffer.WriteByte('?')
		s.args = append(s.args, expr.val)
	case *Predicate: // 表达式
		// 左侧表达式
		_, lp := expr.left.(*Predicate)
		if lp {
			s.stringBuffer.WriteByte('(')
		}
		if err := s.buildExpression(expr.left); err != nil {
			return err
		}
		if lp {
			s.stringBuffer.WriteByte(')')
		}
		// 链接符
		s.stringBuffer.WriteByte(' ')
		s.stringBuffer.WriteString(string(expr.op))
		s.stringBuffer.WriteByte(' ')
		// 右侧表达式
		_, rp := expr.right.(*Predicate)
		if rp {
			s.stringBuffer.WriteByte('(')
		}
		if err := s.buildExpression(expr.right); err != nil {
			return err
		}
		if rp {
			s.stringBuffer.WriteByte(')')
		}
	}
	return nil
}

// From 指定表名，如果是空字符串，那么将会使用默认表名
func (s *ShardingSelector[T]) From(tbl string) *ShardingSelector[T] {
	s.table = tbl
	return s
}

func (s *ShardingSelector[T]) Where(where ...*Predicate) *ShardingSelector[T] {
	s.where = where
	return s
}

func (s *ShardingSelector[T]) GroupBy(columns ...*Column) *ShardingSelector[T] {
	s.groupBy = columns
	return s
}

func (s *ShardingSelector[T]) Having(having *Predicate) *ShardingSelector[T] {
	s.having = having
	return s
}

func (s *ShardingSelector[T]) OrderBy(columns ...*OrderBy) *ShardingSelector[T] {
	s.orderBy = columns
	return s
}

func (s *ShardingSelector[T]) Limit(limit int) *ShardingSelector[T] {
	s.limit = limit
	return s
}

func (s *ShardingSelector[T]) Offset(offset int) *ShardingSelector[T] {
	s.offset = offset
	return s
}
