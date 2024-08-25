package utils

import (
	"time"
)

func FormatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	formattedTime := t.Format("01/02/2006 03:04:05 PM")
	return formattedTime
}
