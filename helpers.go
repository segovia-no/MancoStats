package main

import (
	"pubgstats/pubgDAL"
)

func GetNamesFromPlayerSlice(players []pubgDAL.Player) []string {
	var names []string
	for _, player := range players {
		names = append(names, player.Name)
	}
	return names
}
