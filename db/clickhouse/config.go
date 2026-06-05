package clickhouse

type Config struct {
	DNS      string `properties:"clickhouse.dns"`
	Database string `properties:"clickhouse.database"`
	User     string `properties:"clickhouse.user"`
	Password string `properties:"clickhouse.password"`
}
