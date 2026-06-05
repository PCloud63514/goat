package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	sqldriver "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

type DB interface {
	DB() *sql.DB
	Close() error
}

type DBImpl struct {
	db            *sql.DB
	sshClient     *ssh.Client
	stopKeepalive chan struct{}
}

func NewDB(config *Config) DB {
	dsnConfig := sqldriver.NewConfig()
	dsnConfig.User = config.User
	dsnConfig.Passwd = config.Password
	dsnConfig.DBName = config.Database
	dsnConfig.ParseTime = true
	dsnConfig.Timeout = 10 * time.Second
	dsnConfig.ReadTimeout = 30 * time.Second
	dsnConfig.WriteTimeout = 30 * time.Second
	dsnConfig.Loc = resolveMySQLLocation()

	impl := &DBImpl{}

	if config.SSHHost != "" {
		sshClient, err := newSSHClient(config)
		if err != nil {
			panic(fmt.Errorf("failed to connect via ssh: %w", err))
		}
		impl.sshClient = sshClient

		stop := make(chan struct{})
		impl.stopKeepalive = stop
		go sshKeepalive(sshClient, 30*time.Second, stop)

		dialContextKey := fmt.Sprintf("mysql+ssh:%s", config.SSHHost)
		sqldriver.RegisterDialContext(dialContextKey, func(ctx context.Context, addr string) (net.Conn, error) {
			return sshClient.Dial("tcp", addr)
		})
		dsnConfig.Net = dialContextKey
	} else {
		dsnConfig.Net = "tcp"
	}

	dsnConfig.Addr = fmt.Sprintf("%s:%d", config.Host, config.Port)
	dsn := dsnConfig.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to mysql: %w", err))
	}

	maxOpen := config.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 10
	}
	maxIdle := config.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 5
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetConnMaxIdleTime(10 * time.Second)
	impl.db = db
	return impl
}

func resolveMySQLLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Seoul")
	if err == nil {
		return location
	}
	return time.Local
}

func newSSHClient(config *Config) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: config.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", config.SSHHost, config.SSHPort)
	return ssh.Dial("tcp", addr, sshConfig)
}

func sshKeepalive(client *ssh.Client, interval time.Duration, stop <-chan struct{}) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			client.SendRequest("keepalive@openssh.com", true, nil)
		case <-stop:
			return
		}
	}
}

func (d *DBImpl) DB() *sql.DB {
	return d.db
}

func (d *DBImpl) Close() error {
	if d.stopKeepalive != nil {
		close(d.stopKeepalive)
	}
	err := d.db.Close()
	if d.sshClient != nil {
		d.sshClient.Close()
	}
	return err
}
