package commonerr

import "errors"

var (
	ErrNotFound  = errors.New("not found error")
	ErrForbidden = errors.New("permission denied error")
)
