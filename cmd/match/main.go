package main

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/matcher/glicko2"
	"github.com/hedon954/go-matcher/internal/service/impl"
)

func main() {
	var (
		playerLimit  = 5
		matchChannel = make(chan entry.Group, 1024)
	)

	glicko2 := glicko2.New()
	glicko2 = glicko2
	impl := impl.NewDefault(playerLimit, matchChannel)
	impl = impl
}
