package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const CURRENT_SEASON_ID = "division.bro.official.pc-2018-33"

type GameMode string
type StatsType string

const (
	Squad  GameMode  = "squad"
	Duo    GameMode  = "duo"
	All    StatsType = "all"
	Weekly StatsType = "weekly"
)

func SendHelpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := fmt.Sprintf("Comandos disponibles: \n  season \n  semana \n  vicio \n  playerlist \n  saveplayer [nombre] \n")
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		fmt.Println(err)
	}
}

type StatsCommandFlags struct {
	Season    string
	Mode      GameMode
	StatsType StatsType
}

func SendSeasonStats(players PlayerList, commandTail []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	flags := StatsCommandFlags{ // defaults
		Season:    CURRENT_SEASON_ID,
		Mode:      Squad,
		StatsType: All,
	}

	for i := 0; i < len(commandTail); i++ {
		switch commandTail[i] {
		case "squad":
			flags.Mode = Squad
		case "duo":
			flags.Mode = Duo
		case "all":
			flags.StatsType = All
		case "weekly", "semana":
			flags.StatsType = Weekly
		}
	}

	playerListStats, err := MultiplePlayerStats(players, CURRENT_SEASON_ID, flags.Mode)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Un error ha ocurrido con la API de PUBG")
		return
	}

	var embed *discordgo.MessageEmbed
	switch flags.StatsType {
	case All:
		embed, err = GenPlayerStatsEmbedDiscordMsg(
			playerListStats,
			players,
			EmbedProps{
				Title: fmt.Sprintf("Stats season %v", flags.Mode),
				Color: 0xFF9900,
			},
			generalStatsColumns,
		)
	case Weekly:
		embed, err = GenPlayerStatsEmbedDiscordMsg(
			playerListStats,
			players,
			EmbedProps{
				Title: fmt.Sprintf("Stats semanal %v", flags.Mode),
				Color: 0x1E81B0,
			},
			weeklyStatsColumns,
		)
	}
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

func SendAddictionStats(players PlayerList, s *discordgo.Session, m *discordgo.MessageCreate) {
	playerListStats, err := MultiplePlayerStats(players, CURRENT_SEASON_ID, Squad)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Un error ha ocurrido con la API de PUBG")
		return
	}

	embed, err := GenPlayerStatsEmbedDiscordMsg(
		playerListStats,
		players,
		EmbedProps{
			Title: fmt.Sprintf("Stats vicio season %v", Squad),
			Color: 0xA8233D,
		},
		addictionStatsColumns,
	)
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

func SavePlayer(commandTail []string, serverPlayerList *ServerPlayerList, s *discordgo.Session, m *discordgo.MessageCreate) {
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

	err = serverPlayerList.addPlayer(Player{
		ID:   playerId,
		Name: playerName,
	})
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "No se pudo guardar este jugador, ya estaba guardado o la lista esta llena")
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Jugador guardado: "+playerName)
	if err != nil {
		fmt.Println(err)
	}
}

func RemovePlayer(commandTail []string, serverPlayerList *ServerPlayerList, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(commandTail) != 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Cantidad de argumentos incorrecta")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	playerName := commandTail[1]
	if playerName == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "No se puede eliminar un jugador sin el nombre")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	err := serverPlayerList.removePlayer(Player{
		Name: playerName,
	})
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "No se pudo eliminar este jugador, no estaba guardado")
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Jugador eliminado: "+playerName)
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
