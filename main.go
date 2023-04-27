package main

//static variable for calculation
const (
	K     float64 = 20
	Ratio float64 = 400

	//score definition
	ScoreWin  float64 = 1
	ScoreLose float64 = 0
	ScoreDraw float64 = 0.5

	ScoreWinString  = "win"
	ScoreLoseString = "lose"
	ScoreDrawString = "draw"

	DefaultElo float64 = 800
)

var (
	ScoreDefinitionMap = map[string]float64{
		ScoreWinString:  ScoreWin,
		ScoreLoseString: ScoreLose,
		ScoreDrawString: ScoreDraw,
	}
)

func main() {
	/*
		functionality:
		# add player and elo
		# update elo of a player -> for emergency only
		# delete player
		# fetch elo of all player
		# fetch elo of a player
		# calculate elo of 2 players
	*/

	//calculation is based on this https://www.omnicalculator.com/sports/elo

	//first init the latest data from csv file
	InitRecord()

	//then start the app
	InitRouter()
}
