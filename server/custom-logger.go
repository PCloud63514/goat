package server

import (
	"fmt"
	gin "github.com/gin-gonic/gin"
	"time"
)

type CustomLogFormatterParams struct {
	*gin.LogFormatterParams
}

type LOG_LEVEL string

const (
	LOG_INFO  LOG_LEVEL = "\033[36mINFO\033[0m"
	LOG_ERROR LOG_LEVEL = "\033[31mERRO\033[0m"
)

func CustomLogger() gin.HandlerFunc {
	return CustomLoggerWithConfig(gin.LoggerConfig{})
}

var customLogFormatter = func(param CustomLogFormatterParams) string {
	var statusColor, methodColor, resetColor string

	statusColor = param.StatusCodeColor()
	methodColor = param.MethodColor()
	resetColor = param.ResetColor()

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	var l = LOG_INFO
	if 0 < len(param.ErrorMessage) || 400 <= param.StatusCode {
		l = LOG_ERROR
	}

	return fmt.Sprintf("%v[%v] %13v | %s %3d %s %s %-7s %s %#v\n%s",
		l,
		param.TimeStamp.Format("15:04:05"),
		param.Latency,
		statusColor, param.StatusCode, resetColor,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func CustomLoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
	formatter := customLogFormatter

	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	notLogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		if _, ok := skip[path]; !ok {
			param := CustomLogFormatterParams{
				LogFormatterParams: &gin.LogFormatterParams{
					Request: c.Request,
					Keys:    c.Keys,
				},
			}
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			fmt.Fprint(out, formatter(param))
		}
	}
}
