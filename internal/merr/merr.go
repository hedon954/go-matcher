package merr

import (
	"errors"
)

var (
	ErrGroupFull                   = errors.New("group already full")
	ErrGroupDissolved              = errors.New("group dissolved")
	ErrGroupNotExists              = errors.New("group not exists")
	ErrOnlyCaptainCanDissolveGroup = errors.New("only captain can dissolve group")
	ErrOnlyCaptainCanKickPlayer    = errors.New("only captain can kick player")
	ErrPlayerNotExists             = errors.New("player not exists")
	ErrRoomNotExists               = errors.New("room not exists")
	ErrPlayerNotInRoom             = errors.New("player not in room")
	ErrKickSelf                    = errors.New("cannot kick self")
	ErrChangeSelfRole              = errors.New("cannot change self role")
	ErrNotCaptain                  = errors.New("you not captain")
	ErrPermissionDeny              = errors.New("permission deny")
	ErrInvitationExpired           = errors.New("invitation expired")
	ErrGameModeNotMatch            = errors.New("game mode not match")
	ErrGroupVersionTooLow          = errors.New("group version too low")
	ErrPlayerVersionTooLow         = errors.New("player version too low")

	ErrPlayerOffline    = errors.New("player offline")
	ErrPlayerNotInGroup = errors.New("player not in group")
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
