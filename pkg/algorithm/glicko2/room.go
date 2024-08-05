package glicko2

// Room is an abstract representation of a room, composed of multiple teams
type Room interface {
	// Get the ID of the room
	GetID() int64

	// Get the teams in the room
	GetTeams() []Team

	// Add a team to the room
	AddTeam(t Team)

	// Remove a team from the room
	RemoveTeam(t Team)

	// Get the MMR (Match Market Rating) value of the room
	GetMMR() float64

	// Get the number of players in the room
	PlayerCount() int

	// Get the start time of the match, which is the earliest start time of the player
	GetStartMatchTimeSec() int64

	// Get the end time of the match
	GetFinishMatchTimeSec() int64
	SetFinishMatchTimeSec(t int64)

	// Check if the room contains an AI
	HasAi() bool
}
