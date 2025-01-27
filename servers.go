package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type ServerPlayerList struct {
	ServerID   string
	PlayerList PlayerList
}

type PlayerList []Player

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const ServersPath = "servers"
const CsvSuffix = "_players.csv"

func GetServerIndex(serverID string, serversList *[]ServerPlayerList) (int, error) {
	for i, server := range *serversList {
		if server.ServerID == serverID {
			return i, nil
		}
	}

	// If the server doesn't exist, create it and return the index
	newSrvPtr, err := CreateServer(serverID, serversList)
	if err != nil {
		fmt.Println("[WARN] " + err.Error())
		return 0, err
	}

	return newSrvPtr, nil
}

func CreateServer(serverID string, serversList *[]ServerPlayerList) (int, error) {
	var newServer = ServerPlayerList{
		ServerID:   serverID,
		PlayerList: PlayerList{},
	}

	err := overwritePlayersCSV(serverID+CsvSuffix, newServer.PlayerList)
	if err != nil {
		return 0, err
	}

	*serversList = append(*serversList, newServer)
	return len(*serversList) - 1, nil
}

func LoadServerCSVs(svrPlayerList *[]ServerPlayerList) error {
	var serverFiles []string
	err := filepath.WalkDir(ServersPath, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".csv" {
			serverFiles = append(serverFiles, s)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Load players from each CSV into the provided ServerPlayerList
	for _, serverFile := range serverFiles {
		serverID := filepath.Base(serverFile)
		serverID = serverID[:len(serverID)-len(CsvSuffix)]
		sp := ServerPlayerList{
			ServerID: serverID,
		}
		if err := sp.loadPlayers(); err != nil {
			fmt.Println(err)
			continue
		}
		*svrPlayerList = append(*svrPlayerList, sp)
	}

	return nil
}

func (sp *ServerPlayerList) loadPlayers() error {
	if sp.ServerID == "" {
		return errors.New("server id not set")
	}

	players, err := readPlayersCSV(sp.ServerID + CsvSuffix)
	if err != nil {
		return err
	}
	sp.PlayerList = players
	return nil
}

func (sp *ServerPlayerList) addPlayer(player Player) error {
	if player.Name == "" {
		return errors.New("player name not set")
	}

	for _, existingPlayer := range sp.PlayerList {
		if existingPlayer.ID == player.ID {
			return errors.New("player already exists, skipping add")
		}
	}

	sp.PlayerList = append(sp.PlayerList, player)

	err := overwritePlayersCSV(sp.ServerID+CsvSuffix, sp.PlayerList)
	if err != nil {
		return err
	}
	return nil
}

func (sp *ServerPlayerList) removePlayer(player Player) error {
	for i, existingPlayer := range sp.PlayerList {
		if player.Name == existingPlayer.Name {
			sp.PlayerList = append((sp.PlayerList)[:i], (sp.PlayerList)[i+1:]...)
			err := overwritePlayersCSV(sp.ServerID+CsvSuffix, sp.PlayerList)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("player not found")
}

func readPlayersCSV(filename string) ([]Player, error) {
	fullFilePath := filepath.Join(ServersPath, filename)

	file, err := os.Open(fullFilePath)
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

func overwritePlayersCSV(filename string, players []Player) error {
	err := os.MkdirAll(ServersPath, os.ModePerm)
	if err != nil {
		return err
	}

	fullFilePath := filepath.Join(ServersPath, filename)

	file, err := os.Create(fullFilePath)
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
