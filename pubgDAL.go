package main

import (
	"errors"
	"fmt"
	"strings"
)

func MultiplePlayerStats(players []Player, seasonId string, mode GameMode) ([]PlayerStats, error) {
	fmt.Println("Requesting PUBG stats")
	PubgPlayerIds := GetIdsFromPlayerSlice(players)

	reqUrl := PUBG_API_URL + "/seasons/" + seasonId + "/gameMode/" + string(mode) + "/players?filter[playerIds]=" + strings.Join(PubgPlayerIds, ",")

	var respStats StatsResponse
	err := PubgApiGET(reqUrl, &respStats)
	if err != nil {
		return []PlayerStats{}, err
	}

	if len(respStats.PlayerStatsList) < 1 {
		return []PlayerStats{}, errors.New("no info returned")
	}

	return respStats.PlayerStatsList, nil
}

func FindPlayerIdFromName(name string) (string, error) {
	fmt.Println("Looking for player: " + name)

	reqUrl := PUBG_API_URL + "/players?filter[playerNames]=" + name

	var respPlayer PlayerResponse
	err := PubgApiGET(reqUrl, &respPlayer)
	if err != nil {
		return "", err
	}

	if len(respPlayer.Data) < 1 {
		return "", errors.New("player not found")
	}

	return respPlayer.Data[0].Id, nil
}
