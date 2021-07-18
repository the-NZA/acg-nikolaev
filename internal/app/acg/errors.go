package acg

import "errors"

var (
	ErrNoRequestParams = errors.New("You need to specify required query params")
	ErrNoCategory      = errors.New("Category does not exist yet")
)
