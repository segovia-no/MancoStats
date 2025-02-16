package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"pubgstats/pubgDAL"
	"strings"
)

type StatsType string

const (
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
	Mode      pubgDAL.GameMode
	StatsType StatsType
}

func SendSeasonStats(players []pubgDAL.Player, commandTail []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(players) < 1 {
		fmt.Println("[WARN] No players found for this server")
		s.ChannelMessageSend(m.ChannelID, "Debes guardar jugadores primero para usar este comando")
		return
	}

	flags := StatsCommandFlags{ // defaults
		Season:    PubgDAL.GetSeasonID(),
		Mode:      pubgDAL.Squad,
		StatsType: All,
	}

	for i := 0; i < len(commandTail); i++ {
		switch commandTail[i] {
		case "squad":
			flags.Mode = pubgDAL.Squad
		case "duo":
			flags.Mode = pubgDAL.Duo
		case "all":
			flags.StatsType = All
		case "weekly", "semana":
			flags.StatsType = Weekly
		case "lifetime", "total":
			flags.Season = "lifetime"
		}
	}

	playersStats, err := PubgDAL.MultiplePlayerStats(players, flags.Season == "lifetime", flags.Mode)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Un error ha ocurrido con la API de PUBG")
		return
	}

	var embed *discordgo.MessageEmbed
	switch flags.StatsType {
	case All:
		embed, err = GenPlayerStatsEmbedDiscordMsg(
			playersStats,
			EmbedProps{
				Title: fmt.Sprintf("Stats season %v", flags.Mode),
				Color: 0xFF9900,
			},
			generalStatsColumns,
		)
	case Weekly:
		embed, err = GenPlayerStatsEmbedDiscordMsg(
			playersStats,
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

func SendAddictionStats(players []pubgDAL.Player, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(players) < 1 {
		fmt.Println("[WARN] No players found for this server")
		s.ChannelMessageSend(m.ChannelID, "Debes guardar jugadores primero para usar este comando")
		return
	}

	playersStats, err := PubgDAL.MultiplePlayerStats(players, false, pubgDAL.Squad)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Un error ha ocurrido con la API de PUBG")
		return
	}

	embed, err := GenPlayerStatsEmbedDiscordMsg(
		playersStats,
		EmbedProps{
			Title: fmt.Sprintf("Stats vicio season %v", pubgDAL.Squad),
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

func SendSavedPlayers(playerList []pubgDAL.Player, s *discordgo.Session, m *discordgo.MessageCreate) {
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

	playerId, err := PubgDAL.FindPlayerIdFromName(playerName)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "No se pudo encontrar el jugador "+playerName)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	err = serverPlayerList.addPlayer(pubgDAL.Player{
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

	err := serverPlayerList.removePlayer(pubgDAL.Player{
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
