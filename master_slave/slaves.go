package master_slave

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// 由于不同数据库的DSN不同，因此不同的数据库实现自己的DSN
type dsnResolver interface {
	// ResolveDSN 将DSN进行解析
	ResolveDSN(dsn string) error
	// GetDomain 获取域名
	GetDomain() string
	// ReplaceDomainByIP 将域名替换成IP
	ReplaceDomainByIP(ip string) string
}

type dnsResolver interface {
	// LookupHost 返回域名下对应的IP
	LookupHost(ctx context.Context, domain string) ([]string, error)
}

type SlavesOption func(s *slaves)

var _ dnsResolver = (*net.Resolver)(nil)

type slaveNode struct {
	name string
	db   *sql.DB
}

type slaves struct {
	domain       string        // 从dsn中解析出的domain
	dnsResolver  dnsResolver   // dns域名解析器，根据dsn中domain查询所有从节点信息，本质是dns的域名解析
	closeChan    chan struct{} // 用于传输数据库的关闭信息
	dsnResolver  dsnResolver   // dsn解析器
	slaveArr     []*slaveNode  // 从节点列表
	slavesDSNArr []string      // 从节点DSN列表
	interval     time.Duration // 查询从库状态的心跳时间
	timeout      time.Duration // 使用domain解析从节点的时间
	once         sync.Once     // 关闭数据库仅需要一次
	driver       string        // 数据库驱动
	lock         sync.RWMutex  // 读写slaveArr需要加锁
	idx          uint32        // 访问slaveArr的下标，循环计数
}

func WithDriver(driver string) SlavesOption {
	return func(s *slaves) {
		s.driver = driver
	}
}

func WithNetResolver(netResolver dnsResolver) SlavesOption {
	return func(s *slaves) {
		s.dnsResolver = netResolver
	}
}

func WithDSNResolver(resolver dsnResolver) SlavesOption {
	return func(s *slaves) {
		s.dsnResolver = resolver
	}
}

func WithInterval(timeout time.Duration) SlavesOption {
	return func(s *slaves) {
		s.timeout = timeout
	}
}

func WithTimeout(interval time.Duration) SlavesOption {
	return func(s *slaves) {
		s.interval = interval
	}
}

func NewSlaves(dsn string, options ...SlavesOption) (*slaves, error) {
	s := &slaves{
		dsnResolver: &MysqlDSN{},
		driver:      "mysql",
		closeChan:   make(chan struct{}),
		interval:    time.Second,
		timeout:     time.Second,
	}
	// 执行用户的option
	for _, opt := range options {
		opt(s)
	}
	err := s.dsnResolver.ResolveDSN(dsn)
	if err != nil {
		return nil, err
	}
	// domain
	s.domain = s.dsnResolver.GetDomain()
	// 获取所有的从库列表
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	err = s.getSlaves(ctx)
	cancel()
	if err != nil {
		return nil, err
	}
	// 和所有的从节点维护心跳
	go func() {
		ticker := time.NewTimer(s.interval)
		select {
		case <-ticker.C: // 维护心跳
			ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
			err = s.getSlaves(ctx)
			cancel()
			// 获取列表是 try best
			if err != nil {
				log.Print("get slave arr is fail")
			}
		case <-s.closeChan:
			return
		}
	}()
	return s, nil
}

// 1. dsn解析出domain 2. 一个domain可以对应对个IP地址，将domain解析成对应从库的IP地址列表 3. 将IP地址替换掉域名
func (s *slaves) getSlaves(ctx context.Context) error {
	slavesIP, err := s.dnsResolver.LookupHost(ctx, s.domain)
	if err != nil {
		return err
	}
	slavesArr := make([]*slaveNode, 0)
	slavesDSNArr := make([]string, 0)
	for i, ip := range slavesIP {
		slaveDSN := s.dsnResolver.ReplaceDomainByIP(ip)
		slavesDSNArr = append(slavesDSNArr, slaveDSN)
		db, err := sql.Open(s.driver, slaveDSN)
		if err != nil {
			return err
		}
		slavesArr = append(slavesArr, &slaveNode{
			name: strconv.Itoa(i),
			db:   db,
		})
	}
	// 由于domain相同，可能出现读写冲突
	s.lock.Lock()
	s.slaveArr = slavesArr
	s.slavesDSNArr = slavesDSNArr
	s.lock.Unlock()
	return nil
}

// Next 轮训获取从节点
func (s *slaves) Next() (*slaveNode, error) {
	if len(s.slaveArr) == 0 {
		return nil, errors.New("slave arr is empty")
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	atomic.AddUint32(&s.idx, 1)
	idx := int(s.idx) % len(s.slaveArr)
	return s.slaveArr[idx], nil
}
