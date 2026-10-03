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
	errFailedToCreateItem    = errors.New("failed to create item")
	errFailedToGetAllItems   = errors.New("failed to get all items")
	errFailedToGetItemsByID  = errors.New("failed to get items by id")
	errFailedToUpdateItem    = errors.New("failed to update item")
	errFailedToDeleteItem    = errors.New("failed to delete item")
)
