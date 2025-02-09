package entry

import (
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/log"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type PlayerMgr struct {
	*collection.Manager[string, Player]
}

func NewPlayerMgr() *PlayerMgr {
	return &PlayerMgr{Manager: collection.New[string, Player]()}
}

// Encode encodes all players into a map of game modes to their encoded bytes.
//
//nolint:dupl
func (m *PlayerMgr) Encode() map[constant.GameMode][][]byte {
	res := make(map[constant.GameMode][][]byte, m.Len())
	m.Range(func(id string, p Player) bool {
		bs, err := p.Encode()
		if err != nil {
			log.Error().Any("player", p).Err(err).Msg("failed to encode player")
			return true
		}
		res[p.Base().GameMode] = append(res[p.Base().GameMode], bs)
		return true
	})
	return res
}
