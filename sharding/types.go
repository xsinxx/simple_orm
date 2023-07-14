package sharding

import (
	"github.com/simple_orm/model"
)

type Algorithm interface {
	Sharding(op model.Op, val int64) ([]*DataSource, error)
	Broadcast() ([]*DataSource, error)
}

type Pattern struct {
	Base        int64
	DefaultName string
	IsSharding  bool
}

type DataSource struct {
	DB    string
	Table string
}
