package mysql

type Config struct {
	Host     string `properties:"mysql.host"`
	Port     int    `properties:"mysql.port"`
	Database string `properties:"mysql.database"`
	User     string `properties:"mysql.user"`
	Password string `properties:"mysql.password"`

	SSHHost     string `properties:"mysql.ssh.host" default:""`
	SSHPort     int    `properties:"mysql.ssh.port" default:"0"`
	SSHUser     string `properties:"mysql.ssh.user" default:""`
	SSHPassword string `properties:"mysql.ssh.password" default:""`

	MaxOpenConns int `properties:"mysql.maxOpenConns" default:"10"`
	MaxIdleConns int `properties:"mysql.maxIdleConns" default:"5"`
}
