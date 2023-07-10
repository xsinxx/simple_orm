package simple_orm

type ShardingSelector struct {
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
	tableModel, err := s.r.Get(t)
	if err != nil {
		return nil, err
	}
	s.tableModels = tableModel

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
