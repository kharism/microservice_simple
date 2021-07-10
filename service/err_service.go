package service

const (
	ErrUserNil               = "user is nil or empty"
	ErrUsernameExists        = "username is already taken"
	ErrValidationActionNil   = "validation action cannot be nil or empty"
	ErrValidationActionNA    = "validation action \"%s\" is not available"
	ErrVerifyPasswordInvalid = "invalid username/password"
	ErrCreatedByNil          = "created by is null or empty"
	ErrUpdatedByNil          = "updated by is null or empty"
	ErrWhileFetching         = "error occured while fetching %s (%s) : %s"
)
