package utils

import "errors"

var (
	ErrNameCantEmpty     = errors.New("name can't empty")
	ErrQttyInvalid       = errors.New("qtty invalid")
	ErrPriceInvalid      = errors.New("price invalid")
	ErrCategoryNotExists = errors.New("category not exists")
	ErrDBError           = errors.New("db error")
	ErrDataEmpty         = errors.New("data empty")         // 04
	ErrInvalidRequest    = errors.New("invalid request")    // 14
	ErrFailed            = errors.New("failed")             // 01
	ErrProductNotExists  = errors.New("product not exists") // 40
	ErrInsufficientQtty  = errors.New("insufficient qtty")
)
