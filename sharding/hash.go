package sharding

import (
	"errors"
	"github.com/simple_orm/model"
	"strconv"
)

type Hash struct {
	ShardingKey    string
	ShardingKeySet map[string]struct{}
	DBPattern      *Pattern
	TBPattern      *Pattern
}

func NewHashAlgorithm(shardingKey string, shardingKeySet map[string]struct{}, dBPattern *Pattern, tBPattern *Pattern) Algorithm {
	return &Hash{
		ShardingKey:    shardingKey,
		ShardingKeySet: shardingKeySet,
		DBPattern:      dBPattern,
		TBPattern:      tBPattern,
	}
}

func (h *Hash) Sharding(op model.Op, val int64) ([]*DataSource, error) {
	if h.ShardingKey == "" {
		return []*DataSource{}, errors.New("sharding key is empty")
	}
	// 不是sharding key则广播查找
	if _, ok := h.ShardingKeySet[h.ShardingKey]; !ok {
		return h.Broadcast()
	}
	switch op {
	case model.OpEQ:
		dbName := h.DBPattern.DefaultName
		if h.DBPattern.IsSharding {
			dbName = h.DBPattern.DefaultName + strconv.Itoa(int(val%h.DBPattern.Base))
		}
		tbName := h.TBPattern.DefaultName
		if h.TBPattern.IsSharding {
			tbName = h.TBPattern.DefaultName + strconv.Itoa(int(val%h.TBPattern.Base))
		}
		return []*DataSource{
			&DataSource{
				DB:    dbName,
				Table: tbName,
			},
		}, nil
	case model.OpLT, model.OpGT:
		return h.Broadcast()
	default:
		return []*DataSource{}, nil
	}
}

func (h *Hash) Broadcast() ([]*DataSource, error) {
	if h.DBPattern.IsSharding && h.TBPattern.IsSharding { // 分库分表
		return h.shardingDBAndTable()
	} else if h.DBPattern.IsSharding { // 仅分库
		return h.shardingDB()
	} else if h.TBPattern.IsSharding { // 仅分表
		return h.shardingTable()
	}
	return h.allBroadcast()
}

func (h *Hash) shardingDBAndTable() ([]*DataSource, error) {
	return []*DataSource{
		&DataSource{
			DB:    h.DBPattern.DefaultName,
			Table: h.TBPattern.DefaultName,
		},
	}, nil
}

func (h *Hash) shardingDB() ([]*DataSource, error) {
	res := make([]*DataSource, 0)
	for i := 0; i < int(h.DBPattern.Base); i++ {
		res = append(res, &DataSource{
			DB:    h.DBPattern.DefaultName + strconv.Itoa(i),
			Table: h.TBPattern.DefaultName,
		})
	}
	return res, nil
}

func (h *Hash) shardingTable() ([]*DataSource, error) {
	res := make([]*DataSource, 0)
	for i := 0; i < int(h.TBPattern.Base); i++ {
		res = append(res, &DataSource{
			DB:    h.DBPattern.DefaultName,
			Table: h.TBPattern.DefaultName + strconv.Itoa(i),
		})
	}
	return res, nil
}

func (h *Hash) allBroadcast() ([]*DataSource, error) {
	res := make([]*DataSource, 0)
	for i := 0; i < int(h.DBPattern.Base); i++ {
		dbName := h.DBPattern.DefaultName + strconv.Itoa(i)
		for j := 0; j < int(h.TBPattern.Base); j++ {
			res = append(res, &DataSource{
				DB:    dbName,
				Table: h.TBPattern.DefaultName + strconv.Itoa(j),
			})
		}
	}
	return res, nil
}
