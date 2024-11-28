package domain

type PaginationMeta struct {
	CurrentPage int   `json:"currentPage"`
	TotalPages  int   `json:"totalPages"`
	PageSize    int   `json:"pageSize"`
	TotalCount  int64 `json:"totalCount"`
}
type Response[T any] struct {
	Status     bool            `json:"status"`
	StatusCode int             `json:"statusCode"`
	Message    string          `json:"message"`
	Data       T               `json:"data,omitempty"`
	Meta       *PaginationMeta `json:"meta,omitempty"`
}
