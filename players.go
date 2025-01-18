package main

import (
	"encoding/csv"
	"errors"
	"os"
)

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlayerList []Player

func (p PlayerList) ReadPlayersCSV(filename string) ([]Player, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var players []Player
	for _, record := range records {
		players = append(players, Player{
			ID:   record[0],
			Name: record[1],
		})
	}

	return players, nil
}

func (p *PlayerList) addPlayer(player Player) error {
	for _, existingPlayer := range *p {
		if existingPlayer.ID == player.ID {
			return errors.New("player already exists, skipping add")
		}
	}

	*p = append(*p, player)
	return nil
}

func (p PlayerList) OverwritePlayersCSV(filename string, players []Player) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	for _, player := range players {
		err = csvWriter.Write([]string{player.ID, player.Name})
		if err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return nil
}
