package main

import (
	"math"
)

func CalculateElo(matchResult MatchResult) EloCalculationResult {
	initialRecord1 := RecordMap[matchResult.Player1]
	initialRecord2 := RecordMap[matchResult.Player2]

	if matchResult.Draw {
		return doCalculationDraw(initialRecord1, initialRecord2)
	}

	return doCalculationWin(matchResult.PlayerWin, initialRecord1, initialRecord2)
}

func doCalculationDraw(player1, player2 Record) EloCalculationResult {
	newElo1, newElo2 := coreCalculation(player1.Elo, player2.Elo, ScoreDrawString, ScoreDrawString)
	player1.Elo = newElo1
	player2.Elo = newElo2

	return EloCalculationResult{
		Records: []Record{player1, player2},
	}
}

func doCalculationWin(playerWin string, player1, player2 Record) EloCalculationResult {
	var player1Result, player2Result string

	if playerWin == player1.Player {
		player1Result = ScoreWinString
		player2Result = ScoreLoseString
	} else {
		player1Result = ScoreLoseString
		player2Result = ScoreWinString
	}

	newElo1, newElo2 := coreCalculation(float64(player1.Elo), float64(player2.Elo), player1Result, player2Result)
	player1.Elo = newElo1
	player2.Elo = newElo2

	return EloCalculationResult{
		Records: []Record{player1, player2},
	}
}

func coreCalculation(initialElo1, initialElo2 float64, elo1Result, elo2Result string) (float64, float64) {
	var (
		diff          float64
		ratio         float64
		expectedScore float64
		elo1Score     = ScoreDefinitionMap[elo1Result]
		elo2Score     = ScoreDefinitionMap[elo2Result]
		elo1PlusMinus float64
		elo2PlusMinus float64
		finalDiff     float64
		smallerEloWin bool
	)

	diff = initialElo1 - initialElo2
	ratio = diff / Ratio
	tenPowRatio := math.Pow(10, ratio)
	expectedScore = float64(1) / (float64(1) + tenPowRatio)

	elo1PlusMinus, elo2PlusMinus = decidePlusMinus(initialElo1, initialElo2, elo1Score, elo2Score)

	//check whether player that wins has smaller or bigger elo
	if initialElo1 <= initialElo2 {
		if elo1Score == ScoreWin {
			smallerEloWin = true
		}
	} else {
		if elo2Score == ScoreWin {
			smallerEloWin = true
		}
	}

	diffElo1 := (K * (elo1Score - expectedScore))
	diffElo2 := (K * (elo2Score - expectedScore))
	finalDiff = getFinalDiff(diffElo1, diffElo2, smallerEloWin)

	newElo1 := initialElo1 + (finalDiff * elo1PlusMinus)
	newElo2 := initialElo2 + (finalDiff * elo2PlusMinus)

	return roundFloat(newElo1, 1), roundFloat(newElo2, 1)
}

//increment/decrement of elo1 and elo2 is different.
//if player who has lower elo wins, use the bigger diff as finalDiff
//if player who has lower elo lose, use the smaller diff as finalDiff
func getFinalDiff(diffElo1, diffElo2 float64, smallerEloWin bool) (finalDiff float64) {
	if smallerEloWin {
		if math.Abs(diffElo1) >= math.Abs(diffElo2) {
			finalDiff = roundFloat(math.Abs(diffElo1), 1)
		} else {
			finalDiff = roundFloat(math.Abs(diffElo2), 1)
		}
	} else {
		if math.Abs(diffElo1) <= math.Abs(diffElo2) {
			finalDiff = roundFloat(math.Abs(diffElo1), 1)
		} else {
			finalDiff = roundFloat(math.Abs(diffElo2), 1)
		}
	}

	return finalDiff
}

func decidePlusMinus(initialElo1, initialElo2, elo1Score, elo2Score float64) (elo1PlusMinus, elo2PlusMinus float64) {
	if elo1Score == ScoreWin {
		elo1PlusMinus = 1
		elo2PlusMinus = -1
	} else if elo2Score == ScoreWin {
		elo1PlusMinus = -1
		elo2PlusMinus = 1
	} else {
		if initialElo1 >= initialElo2 {
			elo1PlusMinus = -1
			elo2PlusMinus = 1
		} else {
			elo1PlusMinus = 1
			elo2PlusMinus = -1
		}
	}

	return elo1PlusMinus, elo2PlusMinus
}
