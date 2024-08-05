package glicko2

type GroupState uint8

const (
	GroupStateUnready GroupState = iota // Unready state
	GroupStateQueuing                   // Matching state
	GroupStateMatched                   // Matched state
)

// GroupType is the type of team
type GroupType uint8

const (
	// GroupTypeNotTeam represents no team
	GroupTypeNotTeam GroupType = iota
	// GroupTypeNormalTeam represents a normal team
	GroupTypeNormalTeam
	// GroupTypeUnfriendlyTeam represents an unfriendly team
	GroupTypeUnfriendlyTeam
	// GroupTypeMaliciousTeam represents a malicious team
	GroupTypeMaliciousTeam
)

// Group represents a team,
// players can form teams on their own or a single player will be assigned a team when they start matching,
// the team before and after the match will not be broken up.
type Group interface {

	// GetID returns the team ID
	GetID() string

	// MatchKey returns the unique match queue ID
	MatchKey() string

	// GetPlayers returns the list of players in the team
	GetPlayers() []Player

	// PlayerCount returns the number of players in the team
	PlayerCount() int

	// GetMMR returns the MMR value of the team
	GetMMR() float64

	// GetStar returns the team's rank value
	GetStar() int

	// GetState returns the team's state
	GetState() GroupState
	SetState(state GroupState)

	// GetStartMatchTimeSec returns the start time of the match, which is the earliest start time of the player
	GetStartMatchTimeSec() int64
	SetStartMatchTimeSec(t int64)

	// GetFinishMatchTimeSec returns the end time of the match
	GetFinishMatchTimeSec() int64
	SetFinishMatchTimeSec(t int64)

	// Type returns the team type
	Type() GroupType

	// CanFillAi returns true if the team will be filled with AI and the second return value is the AI team to be filled
	CanFillAi() bool

	// ForceCancelMatch is the logic for handling player cancellation when forced to exit
	ForceCancelMatch(reason string, waitSec int64)

	// IsNewer checks if the team is identified as a newcomer
	IsNewer() bool
}
