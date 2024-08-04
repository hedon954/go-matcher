package goat_game

import "github.com/hedon954/go-matcher/internal/entry"

type Group struct {
	*entry.GroupBase
}

func CreateGroup(base *entry.GroupBase) (entry.Group, error) {
	return &Group{GroupBase: base}, nil
}
