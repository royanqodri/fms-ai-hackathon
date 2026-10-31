package util

import (
	"regexp"
	"strings"

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

func ExtractPlateNumber(text string) string {
	re := regexp.MustCompile(`[A-Z]{1,2}\s\d{1,4}\s[A-Z]{1,3}`)
	return re.FindString(strings.ToUpper(text))
}
