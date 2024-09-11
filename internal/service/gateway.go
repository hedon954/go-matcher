package service

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type PlayerGateway interface {
	Get(uid string) entry.Player
	Remove(uid string)
}

type GroupGateway interface {
	Get(groupID int64) entry.Group
	Remove(groupID int64)
	Create(group pto.CreateGroup) (entry.Group, error)
}

type TeamGateway interface {
	Get(teamID int64) entry.Team
	Remove(teamID int64)
	Create(g entry.Group) (entry.Team, error)
}

type RoomGateway interface {
	Get(roomID int64) entry.Room
	Remove(roomID int64)
	Create(t entry.Team) (entry.Room, error)
}
