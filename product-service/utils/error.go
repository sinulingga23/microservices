package utils

import "errors"

var (
	ErrNameCantEmpty     = errors.New("name can't empty")
	ErrQttyInvalid       = errors.New("qtty invalid")
	ErrPriceInvalid      = errors.New("price invalid")
	ErrCategoryNotExists = errors.New("category not exists")
	ErrDBError           = errors.New("db error")
)
