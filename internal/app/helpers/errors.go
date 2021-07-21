package helpers

import "errors"

var (
	ErrNoRequestParams      = errors.New("You need to specify required query params")
	ErrNoCategory           = errors.New("Category does not exist yet")
	ErrNoPost               = errors.New("Post does not exist yet")
	ErrNoService            = errors.New("Service does not exist yet")
	ErrPostAlreadyExist     = errors.New("Post already with exist")
	ErrServiceAlreadyExist  = errors.New("Service already with exist")
	ErrCategoryAlreadyExist = errors.New("Category already exist")
	ErrInvalidObjectID      = errors.New("ObjectID must be valid")
	ErrEmptyObjectID        = errors.New("You need to specify correct ObjectID")
)
