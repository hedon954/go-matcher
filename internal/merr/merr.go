package merr

import (
	"errors"
)

var (
	ErrCreatePlayer                = errors.New("create player failed")
	ErrGroupFull                   = errors.New("group already full")
	ErrGroupDissolved              = errors.New("group dissolved")
	ErrGroupNotExists              = errors.New("group not exists")
	ErrOnlyCaptainCanDissolveGroup = errors.New("only captain can dissolve group")
	ErrOnlyCaptainCanKickPlayer    = errors.New("only captain can kick player")
	ErrVersionNotMatch             = errors.New("version not match")
	ErrPlayerNotExists             = errors.New("player not exists")
	ErrKickSelf                    = errors.New("cannot kick self")
	ErrChangeSelfRole              = errors.New("cannot change self role")
	ErrNotCaptain                  = errors.New("you not capatin")
	ErrPermissionDeny              = errors.New("permission deny")

	ErrPlayerOffline    = errors.New("player offline")
	ErrPlayerNotInGroup = errors.New("create group first")
	ErrPlayerInGroup    = errors.New("player already in group")
	ErrPlayerInMatch    = errors.New("player matching")
	ErrPlayerInGame     = errors.New("player gaming")
	ErrPlayerInSettle   = errors.New("player settling")

	ErrGroupInInvite = errors.New("group not matching")
	ErrGroupInMatch  = errors.New("group matching")
	ErrGroupInGame   = errors.New("group gaming")

	ErrGroupDenyNearbyJoin = errors.New("group deny nearby join")
	ErrGroupDenyRecentJoin = errors.New("group deny recent join")
)
