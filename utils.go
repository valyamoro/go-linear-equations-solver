package main

import (
	"fmt"
	"strconv"
	"strings"
)

var palette = []string{
	"#FF6B9D", "#4EC9DC", "#A78BFA", "#34D399",
	"#FBBF24", "#F87171", "#60A5FA", "#FB923C",
}

func rgba(hex string, alpha float64) string {
	hex = strings.TrimPrefix(hex, "#")
	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)
	return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", r, g, b, alpha)
}

func prettyEq(s string) string {
	return s
}