package main

import (
	"errors"
	"fmt"
	"strings"
)

func MultiplePlayerStats(players []Player, seasonId string, mode string) (string, error) {
	fmt.Println("Requesting PUBG stats")
	PubgPlayerIds := GetIdsFromPlayerSlice(players)

	reqUrl := PUBG_API_URL + "/seasons/" + seasonId + "/gameMode/" + mode + "/players?filter[playerIds]=" + strings.Join(PubgPlayerIds, ",")

	var respStats StatsResponse
	err := PubgApiGET(reqUrl, &respStats)
	if err != nil {
		return "", err
	}

	var discordResponseMsg string
	for i := 0; i < len(respStats.Data); i++ {
		currPlayerId := respStats.Data[i].Relationships.Player.Data.ID
		playerName, err := FindNameFromId(Players, currPlayerId)
		if err != nil {
			fmt.Println(err)
			continue
		}

		currStats := respStats.Data[i].Attributes.GameModeStats

		kills := currStats.Squad.Kills
		deaths := currStats.Squad.Losses
		kd := float64(kills) / float64(deaths)

		discordResponseMsg += fmt.Sprintf("%s: %f\n", playerName, kd)
	}

	return discordResponseMsg, nil
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
