package utils

import "errors"

var (
	ErrNameCantEmpty     = errors.New("name can't empty")
	ErrQttyInvalid       = errors.New("qtty invalid")
	ErrPriceInvalid      = errors.New("price invalid")
	ErrCategoryNotExists = errors.New("category not exists")
	ErrDBError           = errors.New("db error")
	ErrDataEmpty         = errors.New("data empty")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrFailed            = errors.New("failed")
	ErrProductNotExists  = errors.New("product not exists")
)
