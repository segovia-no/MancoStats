package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const CURRENT_SEASON_ID = "division.bro.official.pc-2018-33"

type GameMode string

const (
	Squad GameMode = "squad"
	Duo   GameMode = "duo"
)

func SendStats(players PlayerList, commandTail []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(commandTail) != 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Cantidad de argumentos incorrecta")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	gameModeStr := strings.ToLower(commandTail[1])
	gameMode := GameMode(gameModeStr)
	if !(gameMode == Duo) && !(gameMode == Squad) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Se debe enviar 'duo' o 'squad'")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	playerListStats, err := MultiplePlayerStats(players, CURRENT_SEASON_ID, gameMode)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Un error ha ocurrido con la API de PUBG")
		return
	}

	embed, err := FormatPlayerStatsAsEmbedDiscordMessage(playerListStats, gameModeStr)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Un error ha ocurrido con el formato de Discord")
		return
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		fmt.Println(err)
	}
}

func SendSavedPlayers(playerList PlayerList, s *discordgo.Session, m *discordgo.MessageCreate) {
	playerNames := GetNamesFromPlayerSlice(playerList)

	playerNamesMsg := fmt.Sprintf("%v Jugadores guardados \n --- \n", len(playerNames))
	playerNamesMsg += strings.Join(playerNames, "\n")

	_, err := s.ChannelMessageSend(m.ChannelID, playerNamesMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func SavePlayer(commandTail []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(commandTail) != 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Cantidad de argumentos incorrecta")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	playerName := commandTail[1]
	if playerName == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "No se puede guardar un jugador sin nombre")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	playerId, err := FindPlayerIdFromName(playerName)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "No se pudo encontrar el jugador "+playerName)
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
		s.ChannelMessageSend(m.ChannelID, "No se pudo guardar este jugador, ya estaba guardado o la lista esta llena")
		return
	}

	err = Players.OverwritePlayersCSV("players.csv", Players)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "No se pudo guardar la lista de jugadores en el CSV")
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Jugador guardado: "+playerName)
	if err != nil {
		fmt.Println(err)
	}
}

func SendUnrecognizedCommandMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Saltando: Comando no reconocido")
	if err != nil {
		fmt.Println(err)
	}
}
