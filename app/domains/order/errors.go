package order

import "errors"

var (
	ErrNoStockAvailable      error = errors.New("no stock available")
	ErrStockUpdateInProgress       = errors.New("stock update in progress")
)
