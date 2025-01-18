package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const CURRENT_SEASON_ID = "division.bro.official.pc-2018-33"

func SendSeasonStats(players PlayerList, s *discordgo.Session, m *discordgo.MessageCreate) {
	msg, err := MultiplePlayerStats(players, CURRENT_SEASON_ID, "squad")
	if msg == "" {
		fmt.Println("No info returned!")
		s.ChannelMessageSend(m.ChannelID, "An error has occurred with the PUBG API")
		return
	}
	_, err = s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		fmt.Println(err)
	}
}

func SendSavedPlayers(playerList PlayerList, s *discordgo.Session, m *discordgo.MessageCreate) {
	playerNames := GetNamesFromPlayerSlice(playerList)

	playerNamesMsg := fmt.Sprintf("%v Saved players \n --- \n", len(playerNames))
	playerNamesMsg += strings.Join(playerNames, "\n")

	_, err := s.ChannelMessageSend(m.ChannelID, playerNamesMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func SavePlayer(playerName string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if playerName == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Cannot save: No player name sent")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	playerId, err := FindPlayerIdFromName(playerName)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Cannot save: Player not found")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	err = Players.addPlayer(Player{
		ID:   playerId,
		Name: playerName,
	})
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Couldn't save players list")
		return
	}

	err = Players.OverwritePlayersCSV("players.csv", Players)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Couldn't save players list")
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Saved player: "+playerName)
	if err != nil {
		fmt.Println(err)
	}
}

func SendUnrecognizedCommandMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Skipping: Unrecognized command")
	if err != nil {
		fmt.Println(err)
	}
}
