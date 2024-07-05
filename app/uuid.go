package app

import "github.com/google/uuid"

type UUIDGenerator func() string

var DefaultUUIDGenerator = func() string {
	return uuid.NewString()
}
