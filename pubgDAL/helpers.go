package pubgDAL

import "errors"

func GetIdsFromPlayerSlice(players []Player) []string {
	var ids []string
	for _, player := range players {
		ids = append(ids, player.ID)
	}
	return ids
}

func SimplifyStatsResponse(players []Player, respStats []PlayerGameModeStats, gameMode GameMode) []PlayerStats {
	var playerStats []PlayerStats
	for _, player := range respStats {
		name, _ := FindNameFromId(players, player.Relationships.Player.Data.ID)

		var stats Stats
		if gameMode == Squad {
			stats = player.Attributes.GameModeStats.Squad
		} else if gameMode == Duo {
			stats = player.Attributes.GameModeStats.Duo
		}

		playerStats = append(playerStats, PlayerStats{
			Name:  name,
			Stats: stats,
		})
	}
	return playerStats
}

func FindNameFromId(players []Player, id string) (string, error) {
	for _, player := range players {
		if player.ID == id {
			return player.Name, nil
		}
	}
	return "", errors.New("player not found")
}
