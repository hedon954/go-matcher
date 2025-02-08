package entry

import (
	"sync/atomic"

	"github.com/hedon954/go-matcher/pkg/collection"
)

type GroupMgr struct {
	*collection.Manager[int64, Group]

	// groupIDIter is a counter used to generate unique group IDs.
	// Before we shut down the server, we need to store the current groupID,
	// and then restart the server from the stored groupID.
	// This is done to avoid conflicts with group IDs that were generated before the server was shut down.
	groupIDIter atomic.Int64
}

// NewGroupMgr creates a group repository, `groupIDStart`: the starting group ID.
func NewGroupMgr(groupIDStart int64) *GroupMgr {
	mgr := &GroupMgr{
		Manager: collection.New[int64, Group](),
	}
	mgr.groupIDIter.Store(groupIDStart)
	return mgr
}

func (m *GroupMgr) GenGroupID() int64 {
	return m.groupIDIter.Add(1)
}
