package character

type CharacterListResponse struct {
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	System string `json:"system"`
}

type CharacterDetailResponse struct {
	Name      string `json:"string"`
	MaxHP     int    `json:"max_hp"`
	CurrentHP int    `json:"current_hp"`
}

type CharacterCurrentHpResponse struct {
	CurrentHP int `json:"current_hp"`
}
