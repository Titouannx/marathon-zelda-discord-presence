package model

type PresenceStatus struct {
	Active           bool   `json:"active"`
	GameID           int    `json:"gameId"`
	GameName         string `json:"gameName"`
	GameLogoURL      string `json:"gameLogoUrl"`
	SessionStartedAt string `json:"sessionStartedAt"`
	ProfileURL       string `json:"profileUrl"`
}
