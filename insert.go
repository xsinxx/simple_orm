package simple_orm

import (
	"context"
	"database/sql"
	"errors"
)

type Insert[T any] struct {
	Builder
	values  []any
	db      *DB
	columns []string
	upsert  *UpsertKey
}

func NewInserter[T any](db *DB) *Insert[T] {
	return &Insert[T]{
		Builder: Builder{
			dialect: db.dialect,
		},
		db: db,
	}
}

func (i *Insert[T]) Values(values ...any) *Insert[T] {
	i.values = values
	return i
}

func (i *Insert[T]) Columns(columns ...string) *Insert[T] {
	i.columns = columns
	return i
}

// ==============================  UpsertBuilder =============================

func (i *Insert[T]) OnDuplicateKey() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		insert: i,
	}
}

func (o *UpsertBuilder[T]) Update(assigns ...Assignable) *Insert[T] {
	o.insert.upsert = &UpsertKey{
		assigns: assigns,
	}
	return o.insert
}

func (i *Insert[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errors.New("[insert] insert zero row")
	}
	var (
		t   T
		err error
	)
	tableModel, err := i.db.r.Get(t)
	if err != nil {
		return nil, err
	}
	i.tableModels = tableModel
	i.sb.WriteString("INSERT INTO")
	// table name
	i.sb.WriteString(" `")
	i.sb.WriteString(i.tableModels.TableName)
	i.sb.WriteString("`")

	// column
	i.sb.WriteString("(")
	if len(i.columns) == 0 { // 若没有传入列就使用所有的列
		i.columns = tableModel.ColumnNames
	}
	for idx, column := range i.columns {
		field, ok := tableModel.Col2Field[column]
		if !ok {
			return nil, errors.New("field not exists")
		}
		i.sb.WriteString("`" + field.ColumnName + "`")
		if idx != len(i.columns)-1 {
			i.sb.WriteString(",")
		}
	}
	i.sb.WriteString(")")

	// values
	i.sb.WriteString(" VALUES")
	for idx, _ := range i.values {
		i.sb.WriteString("(")
		for j, _ := range i.columns { // ?的数量取决于i.columns
			i.sb.WriteString("?")
			if j != len(i.columns)-1 {
				i.sb.WriteString(",")
			}
		}
		i.sb.WriteString(")")
		if idx != len(i.values)-1 {
			i.sb.WriteString(",")
		}
	}

	// args
	for _, val := range i.values {
		internalVal := i.db.creator(val, i.tableModels)
		for _, colName := range i.columns {
			// GetValByColName有两种实现方式反射 & Unsafe，默认是Unsafe
			colVal, err := internalVal.GetValByColName(colName)
			if err != nil {
				return nil, err
			}
			i.args = append(i.args, colVal)
		}
	}

	// upsert
	if i.upsert != nil {
		err = i.dialect.Upsert(&i.Builder, i.upsert)
		if err != nil {
			return nil, err
		}
	}
	i.sb.WriteString(";")
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Insert[T]) Exec(ctx context.Context) (sql.Result, error) {
	query, err := i.Build()
	if err != nil {
		return nil, err
	}
	return i.db.store.Exec(query.SQL, query.Args...)
}
