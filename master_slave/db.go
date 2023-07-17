package master_slave

import (
	"context"
	"database/sql"
)

const (
	master = "master"
)

type masterAndSlaves struct {
	master *sql.DB
	slaves *slaves
}

type ShardingQuery struct {
	SQL  string
	Args []any
}

type masterAndSlavesOptions func(m *masterAndSlaves)

func WithSlaves(slaves *slaves) masterAndSlavesOptions {
	return func(m *masterAndSlaves) {
		m.slaves = slaves
	}
}

func NewMasterSlaves(master *sql.DB, options ...masterAndSlavesOptions) (*masterAndSlaves, error) {
	ms := &masterAndSlaves{
		master: master,
	}
	for _, opt := range options {
		opt(ms)
	}
	return ms, nil
}

func (m *masterAndSlaves) Query(ctx context.Context, query *ShardingQuery) (*sql.Rows, error) {
	var db *sql.DB
	// select支持强制走主节点查询
	_, ok := ctx.Value(master).(bool)
	if ok || m.slaves == nil || len(m.slaves.slaveArr) == 0 {
		db = m.master
	} else {
		slaveNode, err := m.slaves.Next()
		if err != nil {
			return nil, err
		}
		db = slaveNode.db
	}
	return db.QueryContext(ctx, query.SQL, query.Args...)
}
