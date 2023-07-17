package master_slave

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"strings"
)

type MysqlDSN struct {
	cfg    *mysql.Config
	domain string
	port   string
}

func (m *MysqlDSN) ResolveDSN(dsn string) error {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return err
	}
	m.cfg = cfg
	index := strings.Index(cfg.Addr, ":")
	if index == -1 {
		return errors.New("address wrong")
	}
	m.domain = cfg.Addr[:index]
	m.port = cfg.Addr[index+1:]
	return nil
}

func (m *MysqlDSN) GetDomain() string {
	return m.domain
}

func (m *MysqlDSN) ReplaceDomainByIP(ip string) string {
	m.cfg.Addr = ip + ":" + m.port
	return m.cfg.FormatDSN()
}
