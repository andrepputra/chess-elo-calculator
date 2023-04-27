package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	RecordMap    map[string]Record //key is player, value is elo
	RecordHeader = []string{"player", "elo"}
)

func InitRecord() error {
	RecordMap = make(map[string]Record)

	csvFile, _ := os.Open("promo_chess_elo.csv")
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if line[0] == "player" {
			continue
		}

		player := line[0]
		elo := line[1]
		elofloat, _ := strconv.ParseFloat(elo, 64)

		RecordMap[player] = Record{
			Player: player,
			Elo:    elofloat,
		}
	}

	return nil
}

func AddToFileRecord(records []Record) {
	csvFile, err := os.OpenFile("promo_chess_elo.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("An error encountered when opening promo_chess_elo.csv ::", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	var temp [][]string

	for _, record := range records {
		temp = append(temp, []string{record.Player, fmt.Sprint(record.Elo)})
	}

	err = writer.WriteAll(temp)
	if err != nil {
		fmt.Println("An error encountered when write csv ::", err)
	}

	writer.Flush()
}

func RewriteFileRecord() {
	csvFile, err := os.OpenFile("promo_chess_elo.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("An error encountered when opening promo_chess_elo.csv ::", err)
	}
	defer csvFile.Close()

	//empty the file first
	if err := os.Truncate("promo_chess_elo.csv", 0); err != nil {
		fmt.Println("error when emptying promo_chess_elo.csv ::", err)
		return
	}

	writer := csv.NewWriter(csvFile)

	var temp [][]string
	temp = append(temp, RecordHeader)
	for _, record := range RecordMap {
		temp = append(temp, []string{record.Player, fmt.Sprint(record.Elo)})
	}

	err = writer.WriteAll(temp)
	if err != nil {
		fmt.Println("An error encountered when write csv ::", err)
	}

	writer.Flush()
}
