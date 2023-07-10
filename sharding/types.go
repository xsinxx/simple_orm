package sharding

type Algorithm interface {
	Sharding() ([]*DataSource, error)
}

type Pattern struct {
	Base       int64
	Name       string
	IsSharding bool
}

type DataSource struct {
	Name  string
	DB    string
	Table string
}
