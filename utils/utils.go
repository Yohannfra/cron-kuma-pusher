package utils

import "net/url"

func EncodeParam(s string) string {
	return url.QueryEscape(s)
}

func Ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
