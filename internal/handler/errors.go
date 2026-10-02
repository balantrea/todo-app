package handler

import "errors"

var (
	errInvalidInputBody      = errors.New("invalid input body")
	errEmptyAuthHeader       = errors.New("empty auth header")
	errInvalidAuthHeader     = errors.New("invalid auth header")
	errUserIDIsNotFound      = errors.New("user is not found")
	errUserIsNotFound        = errors.New("user is not found")
	errFailedToAssertUserID  = errors.New("failed to assert user id")
	errInvalidIDParam        = errors.New("invalid id param")
	errFailedToCreateUser    = errors.New("failed to create user")
	errFailedToGenerateToken = errors.New("failed to create token")
	errInvalidToken          = errors.New("invalid token")
)
