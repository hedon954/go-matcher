package entry

import (
	"sync/atomic"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/log"
	"github.com/hedon954/go-matcher/pkg/collection"
)

type RoomMgr struct {
	*collection.Manager[int64, Room]
	roomIDIter atomic.Int64
}

// NewRoomMgr creates a room repository.
func NewRoomMgr(roomIDStart int64) *RoomMgr {
	mgr := &RoomMgr{
		Manager: collection.New[int64, Room](),
	}
	mgr.roomIDIter.Store(roomIDStart)
	return mgr
}

// Encode encodes all rooms into a map of game modes to room data.
//
//nolint:dupl
func (m *RoomMgr) Encode() map[constant.GameMode][][]byte {
	res := make(map[constant.GameMode][][]byte, m.Len())
	m.Range(func(i int64, r Room) bool {
		bs, err := r.Encode()
		if err != nil {
			log.Error().Any("room", r).Err(err).Msg("failed to encode room")
			return true
		}
		res[r.Base().GameMode] = append(res[r.Base().GameMode], bs)
		return true
	})
	return res
}
