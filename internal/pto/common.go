package pto

import (
	"encoding/json"
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"google.golang.org/protobuf/proto"
)

// PlayerInfo defines the common information of a player.
// It is always used to initial a player.
type PlayerInfo struct {
	UID         string            `json:"uid" binding:"required"`
	GameMode    constant.GameMode `json:"game_mode" binding:"required"`
	ModeVersion int64             `json:"mode_version" binding:"required"`

	Glicko2Info *Glicko2Info `json:"glicko2_info"`
}

// GameResult defines the common information of a game result.
type GameResult struct {
	RoomID         int64
	StartTime      int64
	EndTime        int64
	ScoreSign      string
	GameMode       constant.GameMode
	MatchStrategy  constant.MatchStrategy
	ModeVersion    int64
	PlayerState    map[string]GamePlayerState
	AIPlayer       map[string]bool
	PlayerMetaInfo map[string]PlayerMetaInfo

	// Result is the game result detail info.
	// It is sent from client and different game mode would be different.
	Result []byte
}

type PlayerMetaInfo struct{}

type GamePlayerState uint8

const (
	Offline GamePlayerState = 0
	Online  GamePlayerState = 1
)

func FromGameResultJson[T any](r *GameResult) (*T, error) {
	var t T
	err := json.Unmarshal(r.Result, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func MustFromGameResultJson[T any](r *GameResult) *T {
	var t T
	err := json.Unmarshal(r.Result, &t)
	if err != nil {
		panic(fmt.Sprintf("unmarshal game json result to *T(type=%T) error: %v", t, err))
	}
	return &t
}

func FromGameResultPb[T any](r *GameResult) (*T, error) {
	var t T
	msg, ok := any(&t).(proto.Message)
	if !ok {
		return nil, fmt.Errorf("type *%T does not implement proto.Message", t)
	}
	err := proto.Unmarshal(r.Result, msg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal game protobuf result to *T(type=%T) error: %v", t, err)
	}
	return &t, nil
}

func MustFromGameResultPb[T any](r *GameResult) *T {
	var t T
	msg, ok := any(&t).(proto.Message)
	if !ok {
		panic(fmt.Sprintf("type *%T does not implement proto.Message", t))
	}
	err := proto.Unmarshal(r.Result, msg)
	if err != nil {
		panic(fmt.Sprintf("unmarshal game protobuf result to *T(type=%T) error: %v", t, err))
	}
	return &t
}
