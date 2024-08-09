package pto

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// EnterGroupSourceType is the source type of entering a group.
type EnterGroupSourceType int

const (
	EnterGroupSourceTypeInvite       EnterGroupSourceType = 0 // invited by other
	EnterGroupSourceTypeNearby       EnterGroupSourceType = 1 // from recent list
	EnterGroupSourceTypeRecent       EnterGroupSourceType = 2 // from nearby list
	EnterGroupSourceTypeFriend       EnterGroupSourceType = 3 // from friend list
	EnterGroupSourceTypeWorldChannel EnterGroupSourceType = 4 // from world channel
	EnterGroupSourceTypeClanChannel  EnterGroupSourceType = 5 // from clan channel
	EnterGroupSourceTypeShare        EnterGroupSourceType = 6 // from share link
)

// EnterGroup is the parameter of entering a group
type EnterGroup struct {
	PlayerInfo
	Source EnterGroupSourceType
}

// CreateGroup is the parameter of creating a group
type CreateGroup struct {
	PlayerInfo
}

// UploadPlayerAttr is the parameter for uploading player attributes needed by game presentation.
type UploadPlayerAttr struct {
	// Attribute is the common information of a player for game presentation.
	Attribute

	// Extra is the extra information of a player needed by different game mode.
	// Here, if you want to do each game mode is independent,
	// you need to use 1+n interfaces (uploadCommonAttr +n * uploadxxxGameAttr),
	// the development efficiency is relatively low.
	//
	// After weighing, it was decided to use a common interface for processing,
	// and then use Extra extension fields for different game modes,
	// in the specific game mode implementation,
	// need to parse and carry out the corresponding processing logic.
	Extra []byte
}

// Attribute is the common information of a player for game display.
// You should define them according to your requirement.
type Attribute struct {
	Nickname string
	Avatar   string
	Star     int64
}

func FromAttrJson[T any](a *UploadPlayerAttr) (*T, error) {
	var t T
	if err := json.Unmarshal(a.Extra, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func MustFromAttrJson[T any](a *UploadPlayerAttr) *T {
	var t T
	err := json.Unmarshal(a.Extra, &t)
	if err != nil {
		panic(fmt.Sprintf("unmarshal attribute json extra to *T(type=%T) error: %v", t, err))
	}
	return &t
}

func FromAttrPb[T any](a *UploadPlayerAttr) (*T, error) {
	var t T
	msg, ok := any(&t).(proto.Message)
	if !ok {
		return nil, fmt.Errorf("type %T does not implement proto.Message", t)
	}
	err := proto.Unmarshal(a.Extra, msg)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func MustFromAttrPb[T any](a *UploadPlayerAttr) *T {
	var t T
	msg, ok := any(&t).(proto.Message)
	if !ok {
		panic(fmt.Sprintf("type *%T does not implement proto.Message", t))
	}
	err := proto.Unmarshal(a.Extra, msg)
	if err != nil {
		panic(fmt.Sprintf("unmarshal attribute protobuf extra to *T(type=%T) error: %v", t, err))
	}
	return &t
}
