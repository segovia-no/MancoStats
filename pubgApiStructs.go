package main

type StatsResponse struct {
	Data []struct {
		Attributes struct {
			GameModeStats gameModeStats `json:"gameModeStats"`
		} `json:"attributes"`
		Relationships relationships `json:"relationships"`
	} `json:"data"`
}

type gameModeStats struct {
	Squad struct {
		Assists             int     `json:"assists"`
		Boosts              int     `json:"boosts"`
		DBNOs               int     `json:"dBNOs"` // Number of enemy players knocked
		DailyKills          int     `json:"dailyKills"`
		DamageDealt         float64 `json:"damageDealt"`
		Days                int     `json:"days"`
		DailyWins           int     `json:"dailyWins"`
		HeadshotKills       int     `json:"headshotKills"`
		Heals               int     `json:"heals"`
		Kills               int     `json:"kills"`
		LongestKill         float64 `json:"longestKill"`
		LongestTimeSurvived float64 `json:"longestTimeSurvived"`
		Losses              int     `json:"losses"`
		MaxKillStreaks      int     `json:"maxKillStreaks"`
		MostSurvivalTime    float64 `json:"mostSurvivalTime"`
		Revives             int     `json:"revives"`
		RideDistance        float64 `json:"rideDistance"`
		RoadKills           int     `json:"roadKills"`
		RoundMostKills      int     `json:"roundMostKills"`
		RoundsPlayed        int     `json:"roundsPlayed"`
		Suicides            int     `json:"suicides"`
		SwimDistance        float64 `json:"swimDistance"`
		TeamKills           int     `json:"teamKills"`
		TimeSurvived        float64 `json:"timeSurvived"`
		Top10s              int     `json:"top10s"`
		VehicleDestroys     int     `json:"vehicleDestroys"`
		WalkDistance        float64 `json:"walkDistance"`
		WeaponsAcquired     int     `json:"weaponsAcquired"`
		WeeklyKills         int     `json:"weeklyKills"`
		WeeklyWins          int     `json:"weeklyWins"`
		Wins                int     `json:"wins"`
	} `json:"squad"`
}

type PlayerResponse struct {
	Data []struct {
		Id         string `json:"id"`
		Attributes struct {
			Name string `json:"name"`
		} `json:"attributes"`
	} `json:"data"`
}

type relationships struct {
	Player dataIDAndType `json:"player"`
}

type dataIDAndType struct {
	Data struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"data"`
}
