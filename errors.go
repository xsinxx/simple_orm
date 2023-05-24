package simple_orm

import "errors"

var (
	ColumnsNotMatch  = errors.New("the number of values in dest must be the same as the number of columns in Row")
	ColumnsNotExists = errors.New("column from db not contained in meta")
)
