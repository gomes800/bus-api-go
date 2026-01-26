package util

import (
	"strconv"
	"strings"
)

func Convert(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
