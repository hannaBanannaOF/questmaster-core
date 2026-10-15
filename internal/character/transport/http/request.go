package character

type CreateCharacterRequest struct {
	Name   string `json:"name"`
	System string `json:"system"`
}

type CreatePremadeCharacterRequest struct {
	Name string `json:"name"`
}

type EditCharacterRequest struct {
	Name string `json:"name"`
}
