package servicemock

import (
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/pto"
)

type GameServerDispatch struct{}

func (d *GameServerDispatch) Dispatch(_ constant.GameMode, _ int64) (pto.GameServerInfo, error) {
	return pto.GameServerInfo{
		Host:     "127.0.0.1",
		Port:     8080, //nolint:mnd
		Protocol: constant.KCP,
	}, nil
}
