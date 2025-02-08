package modes

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game"
	"github.com/hedon954/go-matcher/internal/entry/test_game"
)

func Init() {
	defer entry.PrintModes()

	test_game.RegisterFactory()
	goat_game.RegisterFactory()
}
