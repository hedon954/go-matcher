package glicko2

// Team is an abstract representation of a team, composed of 1-n Groups
type Team interface {
	// Get the list of groups
	GetGroups() []Group

	// Add a group to the team
	AddGroup(group Group)

	// Remove a group from the team
	RemoveGroup(groupId string)

	// Get the number of players in the team
	PlayerCount() int

	// Get the MMR (Match Market Rating) value of the team
	GetMMR() float64

	// Get the rank value of the team
	GetStar() int

	// Get the start time of the match, which is the earliest start time of the player
	GetStartMatchTimeSec() int64

	// Get the end time of the match
	GetFinishMatchTimeSec() int64
	SetFinishMatchTimeSec(t int64)

	// Check if the team is an AI team
	IsAi() bool

	// Check if the team can be filled with AI
	CanFillAi() bool

	// Check if the team is full
	IsFull(teamPlayerLimit int) bool

	// Check if the team is considered a newcomer
	IsNewer() bool
}
