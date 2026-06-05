package restful

import (
	"fmt"
	"time"

	"github.com/PCloud63514/goat/logger"
	"github.com/gin-gonic/gin"
)

// GinMiddleware는 Gin 라우터용 요청 로깅 미들웨어입니다.
// TID를 주입하고 요청 완료 시 goat logger로 로그를 기록합니다.
// net/http 스타일의 Middleware()와 달리 ResponseWriter를 감싸지 않아 WebSocket 업그레이드와 호환됩니다.
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		request, tid := WithTID(c.Request)
		c.Request = request
		c.Writer.Header().Set(tidHeaderName, tid)

		c.Next()

		logger.WithContext(c.Request.Context()).Info(
			fmt.Sprintf("%s  status: %d, path: %s %s",
				formatGinDuration(time.Since(startedAt)),
				c.Writer.Status(),
				c.Request.Method,
				c.Request.URL.Path,
			),
		)
	}
}

func formatGinDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.0fμs", float64(d.Microseconds()))
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}
