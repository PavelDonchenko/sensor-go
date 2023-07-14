package utils

import (
	"strconv"
	"time"
)

func ParseUnixToString(unix string) string {
	intUnix, _ := strconv.ParseInt(unix, 10, 64)

	t := time.Unix(intUnix, 0)

	return t.Format("2006-01-02 15:04:05.000000")
}
