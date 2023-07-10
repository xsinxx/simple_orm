package sharding

import (
	"context"
	"github.com/simple_orm/model"
)

type Algorithm interface {
	Sharding(ctx context.Context, op model.Op, val int64) ([]*DataSource, error)
	Broadcast(ctx context.Context) ([]*DataSource, error)
}

type Pattern struct {
	Base        int64
	DefaultName string
	IsSharding  bool
}

type DataSource struct {
	Name  string
	DB    string
	Table string
}
