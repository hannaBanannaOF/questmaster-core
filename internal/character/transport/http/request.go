package character

type CreateCharacterRequest struct {
	Name   string `json:"name"`
	System string `json:"system"`
	Hp     *int   `json:"hp"`
}

type UpdateHPRequest struct {
	NewHP int `json:"new_hp"`
}

type CreatePremadeCharacterRequest struct {
	Name  string `json:"name"`
	MaxHP *int   `json:"max_hp"`
}

type EditCharacterRequest struct {
	Name  *string `json:"name"`
	MaxHP *int    `json:"max_hp"`
}
