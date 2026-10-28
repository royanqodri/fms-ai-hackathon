package util

import (
	"github.com/google/uuid"
)

func ParseUUID(s string) uuid.UUID {
	if s == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func MustParseUUID(s string) uuid.UUID {
	return uuid.MustParse(s)
}
