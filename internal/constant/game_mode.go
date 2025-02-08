package constant

type GameMode int

const (
	GameModeTest     GameMode = -1
	GameModeGoatGame GameMode = 905
)

var GameModeNames = map[GameMode]string{
	GameModeTest:     "test",
	GameModeGoatGame: "goat_game",
}
