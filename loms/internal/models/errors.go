package models

import "github.com/pkg/errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrStocksNotFound = errors.New("stocks not found")
	ErrInternal       = errors.New("unexpected error")
)
