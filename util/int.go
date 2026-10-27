package util

import (
	"strconv"
	"strings"
)

func ParseInt(s string) int64 {
	if strings.TrimSpace(s) == "" {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int64(v)
}

func IntToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func ToString(i int) string {
	return strconv.Itoa(i)
}
