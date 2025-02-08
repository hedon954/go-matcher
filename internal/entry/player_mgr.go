package entry

import (
	"github.com/hedon954/go-matcher/pkg/collection"
)

type PlayerMgr struct {
	*collection.Manager[string, Player]
}

func NewPlayerMgr() *PlayerMgr {
	return &PlayerMgr{Manager: collection.New[string, Player]()}
}
