package glicko2

// Args encapsulates the three core parameters of the Glicko-2 algorithm.
type Args struct {

	// Player rating, which is a direct measure of a player's ability.
	MMR float64 `json:"mmr"`

	// Rating deviation, which is a measure of the accuracy of the rating.
	// If you're a new player or haven't played in a while, your RD will be high, indicating that your true skill may be far from your rating.
	// If you frequently play, your RD will decrease, indicating that your rating is getting closer to your true skill.
	RD float64 `json:"rd"`

	// Volatility; this is a measure of the fluctuation in a player's rating.
	V float64 `json:"v"`
}
