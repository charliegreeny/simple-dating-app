package match

type SwipeInput struct {
	MatchID string `json:"matchId"`
	Swipe   bool   `json:"swipe"`
}

type SwipeOutput struct {
	Matched bool   `json:"matched"`
	MatchID string `json:"matchId"`
}
