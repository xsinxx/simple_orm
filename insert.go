package simple_orm

import (
	"errors"
	"github.com/simple_orm/model"
	"reflect"
	"strings"
)

type Insert[T any] struct {
	values      []any
	db          *DB
	args        []any
	columns     []string
	sb          strings.Builder
	tableModels *model.TableModel
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
	i.sb.WriteString(";")

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
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}
