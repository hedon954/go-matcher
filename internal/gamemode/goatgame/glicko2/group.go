package glicko2

import (
	"matcher/common"
	"matcher/config"
)

type Group struct {
	*common.GroupBase
}

func NewGroup(groupID int64, p common.Player) (*Group, error) {
	const playerLimit = 3
	base := common.NewGroupBase(groupID, p, config.GroupConfig{PlayerLimit: playerLimit})
	// TODO: other op...
	return &Group{GroupBase: base}, nil
}
