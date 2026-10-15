package character

type CreateCharacterRequest struct {
	Name   string `json:"name"`
	System string `json:"system"`
}

type CreatePremadeCharacterRequest struct {
	Name string `json:"name"`
}

type EditCharacterRequest struct {
	Name *string `json:"name"`
	// How the player is called in the character's campaign; only the player, while in a campaign
	PlayerName *string `json:"player_name"`
}
