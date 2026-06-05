package restful

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// TODO  다음 업데이트에 제거 (v1.2.0에 제거 예정)
// Deprecated (2026-04-20)
func ParsePaginationQuery(c *gin.Context) PaginationQuery {
	page := ParseIntQuery(c, "page")
	limit := ParseIntQuery(c, "limit")
	return PaginationQuery{
		Keyword: c.Query("keyword"),
		Page:    page,
		Limit:   limit,
	}
}

// TODO  다음 업데이트에 제거 (v1.2.0에 제거 예정)
// Deprecated (2026-04-20)
func ParseIntQuery(c *gin.Context, key string) int {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return value
}
