package db

import "fmt"

// Page represents pagination parameters.
type Page struct {
	Number int // 1-based page number
	Size   int // items per page
}

// DefaultPage returns default pagination (page 1, 20 items).
func DefaultPage() Page {
	return Page{Number: 1, Size: 20}
}

// Offset returns the SQL OFFSET value.
func (p Page) Offset() int {
	if p.Number < 1 {
		p.Number = 1
	}
	return (p.Number - 1) * p.Size
}

// Limit returns the SQL LIMIT value.
func (p Page) Limit() int {
	if p.Size <= 0 {
		return 20
	}
	if p.Size > 100 {
		return 100
	}
	return p.Size
}

// LimitOffsetClause returns "LIMIT $n OFFSET $m" clause.
func (p Page) LimitOffsetClause(argStart int) (string, []any) {
	clause := fmt.Sprintf("LIMIT $%d OFFSET $%d", argStart, argStart+1)
	return clause, []any{p.Limit(), p.Offset()}
}

// PagedResult wraps a result set with pagination metadata.
type PagedResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewPagedResult creates a PagedResult from items and total count.
func NewPagedResult[T any](items []T, totalCount int64, page Page) PagedResult[T] {
	totalPages := int(totalCount) / page.Limit()
	if int(totalCount)%page.Limit() > 0 {
		totalPages++
	}
	return PagedResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page.Number,
		PageSize:   page.Limit(),
		TotalPages: totalPages,
	}
}
