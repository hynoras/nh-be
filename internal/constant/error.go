package constant

import "errors"

type ErrorDetail error
type ErrorMessage string

const (
	ErrInvalidSortOrder = "Invalid sort order"
)

var (
	ErrSortOrderShouldBeASCOrDESC ErrorDetail = errors.New("sort order should be ASC or DESC")
	ErrInvalidIDFormat            ErrorDetail = errors.New("invalid id format")
)
