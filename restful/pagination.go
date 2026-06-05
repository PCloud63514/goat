package restful

import "math"

const (
	DefaultPageLimit = 20
	DefaultPageRange = 10
)

// PaginationQuery는 모든 페이지네이션 쿼리가 임베딩해야 하는 공통 구조체입니다.
type PaginationQuery struct {
	Keyword string `form:"keyword"`
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
}

// Normalize는 Page와 Limit을 정규화합니다.
func (p *PaginationQuery) Normalize() (page, limit int) {
	page = p.Page
	limit = p.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = DefaultPageLimit
	}
	return page, limit
}

// QueryMeta는 페이지네이션 메타데이터를 표현합니다.
type QueryMeta struct {
	Page         int  `json:"page"`
	Limit        int  `json:"limit"`
	TotalPages   int  `json:"totalPages"`
	TotalCount   int  `json:"totalCount"`
	CurrentCount int  `json:"currentCount"`
	HasNextPage  bool `json:"hasNextPage"`
	HasPrevPage  bool `json:"hasPrevPage"`
	SPage        int  `json:"sPage"`
	EPage        int  `json:"ePage"`
}

// BuildMeta는 QueryMeta를 생성합니다.
func BuildMeta(page, limit, totalCount, currentCount int) QueryMeta {
	totalPages := 0
	if limit > 0 && totalCount > 0 {
		totalPages = int(math.Ceil(float64(totalCount) / float64(limit)))
	}
	hasPrev := page > 1 && totalPages > 0
	hasNext := page < totalPages
	sPage, ePage := PageRange(page, totalPages)

	return QueryMeta{
		Page:         page,
		Limit:        limit,
		TotalPages:   totalPages,
		TotalCount:   totalCount,
		CurrentCount: currentCount,
		HasNextPage:  hasNext,
		HasPrevPage:  hasPrev,
		SPage:        sPage,
		EPage:        ePage,
	}
}

// PageRange는 페이지 범위(sPage, ePage)를 계산합니다.
func PageRange(page, totalPages int) (int, int) {
	if totalPages == 0 {
		return 0, 0
	}
	group := (page - 1) / DefaultPageRange
	start := group*DefaultPageRange + 1
	end := start + DefaultPageRange - 1
	if end > totalPages {
		end = totalPages
	}
	return start, end
}

// Paginate는 제네릭 페이지네이션 헬퍼 함수입니다.
func Paginate[T any](items []T, page, limit int) ([]T, QueryMeta) {
	totalCount := len(items)
	if totalCount == 0 {
		meta := BuildMeta(page, limit, 0, 0)
		return []T{}, meta
	}

	start := (page - 1) * limit
	if start >= totalCount {
		meta := BuildMeta(page, limit, totalCount, 0)
		return []T{}, meta
	}

	end := start + limit
	if end > totalCount {
		end = totalCount
	}

	current := items[start:end]
	meta := BuildMeta(page, limit, totalCount, len(current))
	return current, meta
}
