package clickhouse

import (
	"crypto/tls"
	"github.com/PCloud63514/goat/logger"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type DB struct {
	conn clickhouse.Conn
}

func NewDB(config *Config) *DB {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{config.DNS},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.User,
			Password: config.Password,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "scglab-data-hub", Version: "0.1"},
			},
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialTimeout:     5 * time.Second,
		ConnMaxLifetime: time.Second,
		Debug:           false,
		Debugf: func(format string, v ...interface{}) {
			logger.Debugf(format, v...)
		},
	})
	if err != nil {
		logger.Errorf("[ClickHouse_DB] Open, Failure to connect database: %v", err)
		panic(err)
	}
	return &DB{
		conn: conn,
	}
}

func (d *DB) Conn() clickhouse.Conn {
	return d.conn
}

func (d *DB) Close() error {
	return d.conn.Close()
}
