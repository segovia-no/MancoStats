package main

import (
	"cmp"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"slices"
	"strings"
)

type EmbedProps struct {
	Title string
	Color int
}

// GenPlayerStatsEmbedDiscordMsg Generates a discord embed message type from a list of player stats,
// takes a statsColumns definition that generates the value of each column,
// the first column calculated value is used for sorting
func GenPlayerStatsEmbedDiscordMsg(
	playerListStats []PlayerStats,
	Players []Player,
	embedProps EmbedProps,
	statsColumns StatsColumns,
) (*discordgo.MessageEmbed, error) {

	if len(statsColumns) < 1 {
		return &discordgo.MessageEmbed{}, errors.New("no stats columns defined, cannot generate")
	}

	if statsColumns[0].Float64ValueFunction == nil {
		return &discordgo.MessageEmbed{}, errors.New("first stats column is not sortable (missing Float64ValueFunction), cannot generate")
	}

	slices.SortFunc(playerListStats,
		func(a, b PlayerStats) int {
			if statsColumns[0].AscOrder {
				return cmp.Compare(statsColumns[0].Float64ValueFunction(b), statsColumns[0].Float64ValueFunction(a))
			} else {
				return cmp.Compare(statsColumns[0].Float64ValueFunction(a), statsColumns[0].Float64ValueFunction(b))
			}
		})

	embed := new(discordgo.MessageEmbed)
	embed.Title = embedProps.Title
	embed.Color = embedProps.Color

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

	// generate all other columns according to definition
	for i := 0; i < len(statsColumns); i++ {
		var resultValueSlice []string
		for _, player := range playerListStats {
			if statsColumns[i].Float64ValueFunction != nil {
				resultValue := statsColumns[i].Float64ValueFunction(player)
				resultValueSlice = append(resultValueSlice, fmt.Sprintf(statsColumns[i].Float64StringFormatter, resultValue))
			} else {
				resultValueSlice = append(resultValueSlice, statsColumns[i].StringValueFunction(player))
			}
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   statsColumns[i].ColumnName,
			Value:  strings.Join(resultValueSlice, "\n"),
			Inline: true,
		})
	}

	return embed, nil
}
