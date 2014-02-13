package main

import (
	"strings"
	)

func CommaStringToSlice(s string) []string {
	return strings.Split(s, ",")
}