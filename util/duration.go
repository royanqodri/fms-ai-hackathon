package util

import (
	"fmt"
	"time"
)

func FormatDuration(duration time.Duration) string {
	seconds := int(duration.Seconds()) % 60
	minutes := int(duration.Minutes()) % 60
	hours := int(duration.Hours()) % 24
	days := int(duration.Hours()) / 24

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%d days, ", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%d hours, ", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%d minutes, ", minutes)
	}
	result += fmt.Sprintf("%d seconds", seconds)

	return result
}
