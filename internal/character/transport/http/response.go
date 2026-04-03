package character

type CharacterListResponse struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	System    string `json:"system"`
	CurrentHP *int   `json:"current_hp"`
	MaxHP     *int   `json:"max_hp"`
}

type CharacterDetailResponse struct {
	Id        int    `json:"id"`
	System    string `json:"system"`
	Name      string `json:"name"`
	MaxHP     *int   `json:"max_hp"`
	CurrentHP *int   `json:"current_hp"`
}

type CharacterCurrentHpResponse struct {
	CurrentHP int `json:"current_hp"`
}
