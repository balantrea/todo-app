package repository

import "errors"

var (
	errTitleCannotBeEmpty     = errors.New("title cannot be empty")
	errFailedToInsertListItem = errors.New("failed to insert list item")
)
