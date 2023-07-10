package sharding

import "errors"

type Hash struct {
	ShardingKey string
	DBPattern   *Pattern
	TBPattern   *Pattern
}

func (h *Hash) Sharding() ([]*DataSource, error) {
	if h.ShardingKey == "" {
		return []*DataSource{}, errors.New("sharding key is empty")
	}

}
