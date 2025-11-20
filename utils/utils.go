package utils

import (
	"maps"
	"net/url"
)

func EncodeParam(s string) string {
	return url.QueryEscape(s)
}

func Ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func MergeMaps(m1, m2 map[string]string) map[string]string {
	result := make(map[string]string, len(m1)+len(m2))

	// Copy contents of m1
	maps.Copy(result, m1)

	// Copy contents of m2 (overwrites on conflict)
	maps.Copy(result, m2)

	return result
}
