package service

import (
	"context"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/pto"
)

// GameServerDispatch dispatches game server info base on game mode and mode version.
type GameServerDispatch interface {
	Dispatch(ctx context.Context, gameMode constant.GameMode, modeVersion int64) (info pto.GameServerInfo, err error)
}
