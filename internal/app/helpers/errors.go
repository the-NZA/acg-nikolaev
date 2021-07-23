package helpers

import "errors"

var (
	ErrNoBodyParams            = errors.New("You need to specify required params")
	ErrNoRequestParams         = errors.New("You need to specify required query params")
	ErrNoCategory              = errors.New("Category does not exist yet")
	ErrNoMatCategory           = errors.New("Material category does not exist yet")
	ErrNoPost                  = errors.New("Post does not exist yet")
	ErrNoMaterial              = errors.New("Material does not exist yet")
	ErrNoService               = errors.New("Service does not exist yet")
	ErrPostAlreadyExist        = errors.New("Post already exist")
	ErrUserAlreadyExist        = errors.New("User already exist")
	ErrEmailAlreadyExist       = errors.New("Email already taken")
	ErrMaterialAlreadyExist    = errors.New("Material already exist")
	ErrServiceAlreadyExist     = errors.New("Service already exist")
	ErrCategoryAlreadyExist    = errors.New("Category already exist")
	ErrMatCategoryAlreadyExist = errors.New("Material category already exist")
	ErrInvalidObjectID         = errors.New("ObjectID must be valid")
	ErrEmptyObjectID           = errors.New("You need to specify correct ObjectID")
)
