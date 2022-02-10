package redis

import "errors"

var (
	errType        = errors.New("failed to convert redis reply to []byte")
	errKeyNotFound = errors.New("key does not exist in redis")
	errKeyExists   = errors.New("key already exists in redis")
)
