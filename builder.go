package simple_orm

import (
	"github.com/simple_orm/model"
	"strings"
)

type Builder struct {
	sb          strings.Builder
	tableModels *model.TableModel
	args        []any
}
