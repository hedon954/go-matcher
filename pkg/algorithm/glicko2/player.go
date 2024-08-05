package glicko2

// Player is an abstract representation of a player
type Player interface {

	// Player ID
	GetID() string

	// Is the player an AI?
	IsAi() bool

	// Get the player's MMR value
	GetMMR() float64

	// Get the player's rank value
	GetStar() int

	// Get the time the player started matching
	GetStartMatchTimeSec() int64
	SetStartMatchTimeSec(t int64)

	// Get the time the player finished matching
	GetFinishMatchTimeSec() int64
	SetFinishMatchTimeSec(t int64)

	// Get the player's rank within their team after the match
	GetRank() int
}
