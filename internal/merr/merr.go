package merr

import (
	"errors"
)

var (
	ErrCreatePlayer                = errors.New("Create Player Failed")
	ErrGroupFull                   = errors.New("Group Already Full")
	ErrGroupDissolved              = errors.New("Group Dissolved")
	ErrGroupNotExists              = errors.New("Group Not Exists")
	ErrOnlyCaptainCanDissolveGroup = errors.New("Only Captain Can Dissolve Group")
	ErrVersionNotMatch             = errors.New("Version Not Match")

	ErrPlayerOffline    = errors.New("Player Offline")
	ErrPlayerNotInGroup = errors.New("Create Group First")
	ErrPlayerInGroup    = errors.New("Player Already In Group")
	ErrPlayerInMatch    = errors.New("Player Matching")
	ErrPlayerInGame     = errors.New("Player Gaming")
	ErrPlayerInSettle   = errors.New("Player Settling")

	ErrGroupInInvite = errors.New("Group Not Matching")
	ErrGroupInMatch  = errors.New("Group Matching")
	ErrGroupInGame   = errors.New("Group Gaming")
)
