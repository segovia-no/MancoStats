package main

import (
	"cmp"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"slices"
	"strings"
)

func FormatPlayerStatsAsEmbedDiscordMessage(playerListStats []PlayerStats, Players []Player, gameModeStr string) (*discordgo.MessageEmbed, error) {
	slices.SortFunc(playerListStats,
		func(a, b PlayerStats) int {
			return cmp.Compare(GetKDRatioFromPlayerStats(b), GetKDRatioFromPlayerStats(a))
		})

	embed := new(discordgo.MessageEmbed)
	embed.Title = "Estadísticas de season (" + gameModeStr + ")"
	embed.Color = 0xFF9900

	// player name column
	var playerNameSlice []string
	for i, player := range playerListStats {
		playerName, err := FindNameFromId(Players, player.Relationships.Player.Data.ID)

		switch i {
		case 0:
			playerName = "🥇" + playerName
		case 1:
			playerName = "🥈" + playerName
		case 2:
			playerName = "🥉" + playerName
		}

		if err != nil {
			fmt.Println(err)
			return &discordgo.MessageEmbed{}, err
		}
		playerNameSlice = append(playerNameSlice, playerName)
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Nombre",
		Value:  strings.Join(playerNameSlice, "\n"),
		Inline: true,
	})

	// kd column
	var kdColumnSlice []string
	for _, player := range playerListStats {
		kd := GetKDRatioFromPlayerStats(player)
		kdColumnSlice = append(kdColumnSlice, fmt.Sprintf("%.2f", kd))
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "K/D",
		Value:  strings.Join(kdColumnSlice, "\n"),
		Inline: true,
	})

	// wins/matches column
	var wlColumnSlice []string
	for _, player := range playerListStats {
		wins := player.Attributes.GameModeStats.Squad.Wins
		losses := player.Attributes.GameModeStats.Squad.Losses
		wlColumnSlice = append(wlColumnSlice, fmt.Sprintf("%d/%d", wins, losses))
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "🍗/Partidas",
		Value:  strings.Join(wlColumnSlice, "\n"),
		Inline: true,
	})

	return embed, nil
}
