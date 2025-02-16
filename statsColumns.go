package main

import (
	"fmt"
	"math"
	"pubgstats/pubgDAL"
)

type StatsColumns []StatsColumnDefinition

type StatsColumnDefinition struct {
	ColumnName             string
	Float64ValueFunction   func(pubgDAL.PlayerStats) float64
	StringValueFunction    func(pubgDAL.PlayerStats) string
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

var kdRatioColumn = StatsColumnDefinition{
	ColumnName: "K/D",
	Float64ValueFunction: func(player pubgDAL.PlayerStats) float64 {
		return math.Round(float64(player.Stats.Kills)/float64(player.Stats.Losses)*100) / 100
	},
	AscOrder:               true,
	Float64StringFormatter: "%.2f",
}

var chickenMatchesRatioColumn = StatsColumnDefinition{
	ColumnName: "🍗/Partidas",
	StringValueFunction: func(player pubgDAL.PlayerStats) string {
		wins := player.Stats.Wins
		losses := player.Stats.Losses
		return fmt.Sprintf("%d/%d", wins, losses)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.2f",
}

var weeklyKillsColumn = StatsColumnDefinition{
	ColumnName: "Kills en la semana",
	Float64ValueFunction: func(player pubgDAL.PlayerStats) float64 {
		return float64(player.Stats.WeeklyKills)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.0f",
}

var weeklyWinsColumn = StatsColumnDefinition{
	ColumnName: "🍗 en la semana",
	Float64ValueFunction: func(player pubgDAL.PlayerStats) float64 {
		return float64(player.Stats.WeeklyWins)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.0f",
}

var roundsPlayed = StatsColumnDefinition{
	ColumnName: "Rondas jugadas",
	Float64ValueFunction: func(player pubgDAL.PlayerStats) float64 {
		return float64(player.Stats.RoundsPlayed)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.0f",
}

var daysPlayed = StatsColumnDefinition{
	ColumnName: "Dias jugados",
	Float64ValueFunction: func(player pubgDAL.PlayerStats) float64 {
		return float64(player.Stats.Days)
	},
	AscOrder:               true,
	Float64StringFormatter: "%.0f",
}
