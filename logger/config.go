package logger

type Config struct {
	Level string `properties:"logger.level" default:"info"`
}
