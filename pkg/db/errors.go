package db

import "errors"


var (
	ErrUserNotFound = errors.New("user not found")
	ErrCollectionNotFound = errors.New("collection not found")
	ErrItemNotFound = errors.New("item not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrTeleIdAlreadyExists = errors.New("telegram already registered")
	ErrInvalidUserId   = errors.New("invalid userid")
)