package mysql

type Scanner interface {
	Scan(dest ...interface{}) error
}
