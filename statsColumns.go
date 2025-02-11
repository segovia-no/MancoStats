package main

import (
	"fmt"
	"math"
)

type StatsColumns []StatsColumnDefinition

type StatsColumnDefinition struct {
	ColumnName             string
	Float64ValueFunction   func(PlayerStats) float64
	StringValueFunction    func(PlayerStats) string
	AscOrder               bool
	Float64StringFormatter string
}

var generalStatsColumns = StatsColumns{
	kdRatioColumn,
	chickenMatchesRatioColumn,
}

var weeklyStatsColumns = StatsColumns{
	weeklyKillsColumn,
	weeklyWinsColumn,
}

var addictionStatsColumns = StatsColumns{
	roundsPlayed,
	daysPlayed,
}

// TODO: How to apply squad and duo selections efficiently
var kdRatioColumn = StatsColumnDefinition{
	ColumnName: "K/D",
	Float64ValueFunction: func(player PlayerStats) float64 {
		return math.Round(float64(player.Attributes.GameModeStats.Squad.Kills)/float64(player.Attributes.GameModeStats.Squad.Losses)*100) / 100
	},
	AscOrder:               true,
	Float64StringFormatter: "%.2f",
}

var chickenMatchesRatioColumn = StatsColumnDefinition{
	ColumnName: "🍗/Partidas",
	StringValueFunction: func(player PlayerStats) string {
		wins := player.Attributes.GameModeStats.Squad.Wins
		losses := player.Attributes.GameModeStats.Squad.Losses
		return fmt.Sprintf("%d/%d", wins, losses)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.2f",
}

var weeklyKillsColumn = StatsColumnDefinition{
	ColumnName: "Kills en la semana",
	Float64ValueFunction: func(player PlayerStats) float64 {
		return float64(player.Attributes.GameModeStats.Squad.WeeklyKills)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.0f",
}

var weeklyWinsColumn = StatsColumnDefinition{
	ColumnName: "🍗 en la semana",
	Float64ValueFunction: func(player PlayerStats) float64 {
		return float64(player.Attributes.GameModeStats.Squad.WeeklyWins)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.0f",
}

var roundsPlayed = StatsColumnDefinition{
	ColumnName: "Rondas jugadas",
	Float64ValueFunction: func(player PlayerStats) float64 {
		return float64(player.Attributes.GameModeStats.Squad.RoundsPlayed)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.0f",
}

var daysPlayed = StatsColumnDefinition{
	ColumnName: "Dias jugados",
	Float64ValueFunction: func(player PlayerStats) float64 {
		return float64(player.Attributes.GameModeStats.Squad.Days)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.0f",
}
