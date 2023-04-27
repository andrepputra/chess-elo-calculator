package main

type (
	Record struct {
		Player string  `json:"player"`
		Elo    float64 `json:"elo"`
	}

	MatchResult struct {
		Player1   string `json:"player_1"`
		Player2   string `json:"player_2"`
		Win       bool   `json:"win"`
		Draw      bool   `json:"draw"`
		PlayerWin string `json:"player_win"`
	}

	EloCalculationResult struct {
		Records []Record `json:"records"`
	}
)
