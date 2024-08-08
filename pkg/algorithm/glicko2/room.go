package glicko2

// Room is an abstract representation of a room, composed of multiple teams
type Room interface {
	// Get the teams in the room
	GetTeams() []Team

	// Add a team to the room
	AddTeam(t Team)

	// Get the MMR (Match Market Rating) value of the room
	GetMMR() float64

	// Get the start time of the match, which is the earliest start time of the player
	GetStartMatchTimeSec() int64

	// Check if the room contains an AI
	HasAi() bool
}
