package slack

var slack = New(nil)

func SetSlack(config *Config) {
	slack = New(config)
}

func SetMode(mode string) error {
	if !slack.initialized {
		return ErrNotConfigured
	}
	return (&slack).SetMode(mode)
}

func Mode() string {
	return slack.Mode()
}

func Info(args ...interface{}) error {
	if !slack.initialized {
		return ErrNotConfigured
	}
	return slack.Info(args...)
}

func Infof(format string, args ...interface{}) error {
	if !slack.initialized {
		return ErrNotConfigured
	}
	return slack.Infof(format, args...)
}

func Error(args ...interface{}) error {
	if !slack.initialized {
		return ErrNotConfigured
	}
	return slack.Error(args...)
}

func Errorf(format string, args ...interface{}) error {
	if !slack.initialized {
		return ErrNotConfigured
	}
	return slack.Errorf(format, args...)
}

func WithFields(fields ...AttachmentField) Slack {
	if !slack.initialized {
		return Slack{}
	}
	return slack.WithFields(fields...)
}
