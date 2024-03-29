package simple_orm

import (
	"errors"
	"fmt"
	"github.com/simple_orm/model"
	"github.com/simple_orm/sharding"
	"github.com/valyala/bytebufferpool"
)

type ShardingBuilder struct {
	algorithm    sharding.Algorithm
	stringBuffer *bytebufferpool.ByteBuffer
	tableModels  *model.TableModel
	args         []any
}

func (s *ShardingBuilder) FindDataSource(where ...*Predicate) ([]*sharding.DataSource, error) {
	// 有where语句尝试分库分表
	if len(where) > 0 {
		predicate := where[0]
		for _, v := range where {
			predicate.And(v)
		}
		return s.findDataSourceByAlgorithm(predicate)
	}
	// 无where语句则广播
	return s.algorithm.Broadcast()
}

func (s *ShardingBuilder) findDataSourceByAlgorithm(predicate *Predicate) ([]*sharding.DataSource, error) {
	switch predicate.op {
	case model.OpAnd:
		// 左侧
		left, ok := predicate.left.(*Predicate)
		if !ok {
			return []*sharding.DataSource{}, errors.New("left is not a predicate")
		}
		leftDataSource, err := s.findDataSourceByAlgorithm(left)
		if err != nil {
			return []*sharding.DataSource{}, err
		}
		// 右侧
		right, ok := predicate.right.(*Predicate)
		if !ok {
			return []*sharding.DataSource{}, errors.New("right is not a predicate")
		}
		rightDataSource, err := s.findDataSourceByAlgorithm(right)
		if err != nil {
			return []*sharding.DataSource{}, err
		}
		return intersection(leftDataSource, rightDataSource), nil
	case model.OpOr: // id > 10 or age >= 30
		// 左侧
		left, ok := predicate.left.(*Predicate)
		if !ok {
			return []*sharding.DataSource{}, errors.New("left is not a predicate")
		}
		leftDataSource, err := s.findDataSourceByAlgorithm(left)
		if err != nil {
			return []*sharding.DataSource{}, err
		}
		// 右侧
		right, ok := predicate.right.(*Predicate)
		if !ok {
			return []*sharding.DataSource{}, errors.New("right is not a predicate")
		}
		rightDataSource, err := s.findDataSourceByAlgorithm(right)
		if err != nil {
			return []*sharding.DataSource{}, err
		}
		return union(leftDataSource, rightDataSource), nil
	case model.OpEQ, model.OpGT, model.OpLT:
		left, ok := predicate.left.(*Column)
		if !ok {
			return []*sharding.DataSource{}, errors.New("left is not a column")
		}
		if _, ok = s.tableModels.Col2Field[left.name]; !ok {
			return []*sharding.DataSource{}, errors.New("illegal field")
		}
		right, ok := predicate.left.(*Value)
		if !ok {
			return []*sharding.DataSource{}, errors.New("right is not a value")
		}
		val, ok := right.val.(int)
		if !ok {
			return []*sharding.DataSource{}, errors.New("right is not int type")
		}
		return s.algorithm.Sharding(predicate.op, int64(val))
	default:
		return []*sharding.DataSource{}, nil
	}
}

func getKey(datsSource *sharding.DataSource) string {
	return fmt.Sprintf("DB:%s, Table:%s", datsSource.DB, datsSource.Table)
}

func intersection(left, right []*sharding.DataSource) []*sharding.DataSource {
	res := make([]*sharding.DataSource, 0)
	intersectionMap := map[string]*sharding.DataSource{}
	for _, v := range left {
		key := getKey(v)
		intersectionMap[key] = v
	}
	for _, v := range right {
		key := getKey(v)
		if _, ok := intersectionMap[key]; ok {
			res = append(res, v)
		}
	}
	return res
}

func union(left, right []*sharding.DataSource) []*sharding.DataSource {
	res := make([]*sharding.DataSource, 0)
	unionMap := map[string]*sharding.DataSource{}
	for _, v := range left {
		key := getKey(v)
		unionMap[key] = v
	}
	for _, v := range right {
		key := getKey(v)
		unionMap[key] = v
	}

	for _, v := range unionMap {
		res = append(res, v)
	}

	return res
}
