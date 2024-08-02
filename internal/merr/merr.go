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
	ErrOnlyCaptainCanKickPlayer    = errors.New("Only Captain Can Kick Player")
	ErrVersionNotMatch             = errors.New("Version Not Match")
	ErrPlayerNotExists             = errors.New("Player Not Exists")
	ErrKickSelf                    = errors.New("Cannot Kick Self")
	ErrHandoverSelf                = errors.New("Cannot Handover Captain to Self")
	ErrNotCaptain                  = errors.New("You Not Capatin")
	ErrPermissionDeny              = errors.New("Permission Deny")

	ErrPlayerOffline    = errors.New("Player Offline")
	ErrPlayerNotInGroup = errors.New("Create Group First")
	ErrPlayerInGroup    = errors.New("Player Already In Group")
	ErrPlayerInMatch    = errors.New("Player Matching")
	ErrPlayerInGame     = errors.New("Player Gaming")
	ErrPlayerInSettle   = errors.New("Player Settling")

	ErrGroupInInvite = errors.New("Group Not Matching")
	ErrGroupInMatch  = errors.New("Group Matching")
	ErrGroupInGame   = errors.New("Group Gaming")

	ErrGroupDenyNearbyJoin = errors.New("Group Deny Nearby Join")
	ErrGroupDenyRecentJoin = errors.New("Group Deny Recent Join")
)
