package main

import (
	"errors"
	"math"
	"strings"
)

func MultiplePlayerStats(players []Player, seasonId string, mode GameMode) ([]PlayerStats, error) {
	PubgPlayerIds := GetIdsFromPlayerSlice(players)
	reqAmount := int(math.Ceil(float64(len(PubgPlayerIds)) / 10)) // each request can handle 10 players

	var respStats []PlayerStats
	for i := 0; i < reqAmount; i++ {
		maxIdx := i*10 + 10
		if maxIdx > len(PubgPlayerIds) {
			maxIdx = len(PubgPlayerIds)
		}
		currPlayerListIds := PubgPlayerIds[i*10 : maxIdx]
		reqUrl := PUBGApiURL + "/seasons/" + seasonId + "/gameMode/" + string(mode) + "/players?filter[playerIds]=" + strings.Join(currPlayerListIds, ",")

		var currRespStats StatsResponse
		err := PubgApiGET(reqUrl, &currRespStats)
		if err != nil {
			return []PlayerStats{}, err
		}

		respStats = append(respStats, currRespStats.PlayerStatsList...)
	}

	if len(respStats) < 1 {
		return []PlayerStats{}, errors.New("no info returned")
	}

	return respStats, nil
}

func FindPlayerIdFromName(name string) (string, error) {
	reqUrl := PUBGApiURL + "/players?filter[playerNames]=" + name

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
