// Package main is the entry of matcher
package main

import (
	"github.com/hedon954/go-matcher/config"
	"github.com/hedon954/go-matcher/matcher"
)

func main() {
	matcher.New(nil, config.GroupConfig{})
}
