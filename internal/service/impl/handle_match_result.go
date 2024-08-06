package impl

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/entry"
)

func (impl *Impl) waitForMatchResult() {
	go func() {
		for r := range impl.roomChannel {
			impl.HandleMatchResult(r)
		}
	}()
}

func (impl *Impl) clearDelayTimer(r entry.Room) {
	fmt.Println("handle match result: ", r)
	for _, t := range r.Base().GetTeams() {
		for _, g := range t.Base().GetGroups() {
			fmt.Println("remove timer: ", g.ID())
			impl.removeWaitAttrTimer(g.ID())
			impl.removeWaitAttrTimer(g.ID())
			impl.removeCancelMatchTimer(g.ID())
		}
	}
}
