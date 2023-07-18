package utils

import (
	"regexp"
	"strconv"
)

func ParseCodename(s string) (string, int) {
	r := regexp.MustCompile(`^(\D+)(\d+)$`)
	matches := r.FindStringSubmatch(s)

	alpha := matches[1]
	numStr := matches[2]
	num, _ := strconv.Atoi(numStr)

	return alpha, num
}
