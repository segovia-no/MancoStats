package pubgDAL

import (
	"errors"
	"log"
	"math"
	"strings"
)

type PUBGApiDAL struct {
	apiURL        string
	seasonID      string
	pubgRequester *PUBGRequester
}

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GameMode string

const (
	Squad GameMode = "squad"
	Duo   GameMode = "duo"
)

func NewPUBGApiDAL(apiURL string, apiToken string) (*PUBGApiDAL, error) {
	requester := NewPUBGRequester(apiToken)

	apiDAL := &PUBGApiDAL{
		apiURL:        apiURL,
		seasonID:      "",
		pubgRequester: requester,
	}

	season, err := apiDAL.FindCurrentSeasonId()
	if err != nil {
		return nil, err
	}

	apiDAL.seasonID = season
	log.Println("Current PUBG season ID is: " + season)

	return apiDAL, nil
}

func (p *PUBGApiDAL) GetSeasonID() string {
	return p.seasonID
}

func (p *PUBGApiDAL) MultiplePlayerStats(players []Player, lifetime bool, mode GameMode) ([]PlayerStats, error) {
	PubgPlayerIds := GetIdsFromPlayerSlice(players)
	reqAmount := int(math.Ceil(float64(len(PubgPlayerIds)) / 10)) // each request can handle 10 players

	seasonId := p.seasonID
	if lifetime {
		seasonId = "lifetime"
	}

	var respStats []PlayerGameModeStats
	for i := 0; i < reqAmount; i++ {
		maxIdx := i*10 + 10
		if maxIdx > len(PubgPlayerIds) {
			maxIdx = len(PubgPlayerIds)
		}
		currPlayerListIds := PubgPlayerIds[i*10 : maxIdx]
		reqUrl := p.apiURL + "/seasons/" + seasonId + "/gameMode/" + string(mode) + "/players?filter[playerIds]=" + strings.Join(currPlayerListIds, ",")

		var currRespStats StatsResponse
		err := p.pubgRequester.Get(reqUrl, &currRespStats)
		if err != nil {
			return []PlayerStats{}, err
		}

		respStats = append(respStats, currRespStats.PlayerStatsList...)
	}

	if len(respStats) < 1 {
		return []PlayerStats{}, errors.New("no info returned")
	}

	playerStats := SimplifyStatsResponse(players, respStats, mode)

	return playerStats, nil
}

func (p *PUBGApiDAL) FindPlayerIdFromName(name string) (string, error) {
	reqUrl := p.apiURL + "/players?filter[playerNames]=" + name

	var respPlayer PlayerResponse
	err := p.pubgRequester.Get(reqUrl, &respPlayer)
	if err != nil {
		return "", err
	}

	if len(respPlayer.Data) < 1 {
		return "", errors.New("player not found")
	}

	return respPlayer.Data[0].Id, nil
}

func (p *PUBGApiDAL) FindCurrentSeasonId() (string, error) {
	reqUrl := p.apiURL + "/seasons"

	var respSeasons SeasonsResponse
	err := p.pubgRequester.Get(reqUrl, &respSeasons)
	if err != nil {
		return "", err
	}

	if len(respSeasons.Data) < 1 {
		return "", errors.New("response had no seasons data")
	}

	for _, season := range respSeasons.Data {
		if season.Attributes.IsCurrentSeason {
			return season.ID, nil
		}
	}
	return "", errors.New("season not found")
}
