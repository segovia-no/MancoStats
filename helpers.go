package main

import "errors"

func GetIdsFromPlayerSlice(players []Player) []string {
	var ids []string
	for _, player := range players {
		ids = append(ids, player.ID)
	}
	return ids
}

func GetNamesFromPlayerSlice(players []Player) []string {
	var names []string
	for _, player := range players {
		names = append(names, player.Name)
	}
	return names
}

func FindNameFromId(players []Player, id string) (string, error) {
	for _, player := range players {
		if player.ID == id {
			return player.Name, nil
		}
	}
	return "", errors.New("Player not found")
}
