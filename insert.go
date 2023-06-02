package simple_orm

import (
	"errors"
	"github.com/simple_orm/model"
	"reflect"
	"strings"
)

type Insert[T any] struct {
	values         []any
	db             *DB
	args           []any
	columns        []string
	sb             strings.Builder
	tableModels    *model.TableModel
	onDuplicateKey *OnDuplicateKey
}

func NewInserter[T any](db *DB) *Insert[T] {
	return &Insert[T]{
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

// ==============================  OnDuplicateKeyBuilder =============================

func (i *Insert[T]) OnDuplicateKey() *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
		insert: i,
	}
}

func (o *OnDuplicateKeyBuilder[T]) Update(assigns ...Assignable) *Insert[T] {
	o.insert.onDuplicateKey = &OnDuplicateKey{
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
	columnSet := map[string]struct{}{}
	for _, column := range i.columns {
		columnSet[column] = struct{}{}
	}
	for _, val := range i.values {
		reflectVal := reflect.ValueOf(val).Elem()
		reflectTyp := reflect.TypeOf(val).Elem()
		for j := 0; j < reflectVal.NumField(); j++ {
			fieldName := reflectTyp.Field(j).Name
			if _, ok := columnSet[fieldName]; !ok {
				continue
			}
			i.args = append(i.args, reflectVal.Field(j).Interface())
		}
	}

	// on duplicate key update
	if i.onDuplicateKey != nil {
		i.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
		for _, assign := range i.onDuplicateKey.assigns {
			err = i.BuildExpression(assign)
			if err != nil {
				return nil, err
			}
		}
	}
	i.sb.WriteString(";")
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Insert[T]) BuildExpression(assign Assignable) error {
	switch e := assign.(type) {
	case *Column:
		columnName := e.name
		field, ok := i.tableModels.Col2Field[columnName]
		if !ok {
			return errors.New("column name not exists")
		}
		// 详细结构参照单测
		i.sb.WriteString("`")
		i.sb.WriteString(field.ColumnName)
		i.sb.WriteString("`=VALUES(`")
		i.sb.WriteString(field.ColumnName)
		i.sb.WriteString("`)")
	case *Assignment:
		columnName := e.ColumnName
		field, ok := i.tableModels.Col2Field[columnName]
		if !ok {
			return errors.New("column name not exists")
		}
		i.sb.WriteString("`")
		i.sb.WriteString(field.ColumnName)
		i.sb.WriteString("`")
		i.sb.WriteString("=?")
		i.args = append(i.args, e.Val)
	}
	return nil
}
