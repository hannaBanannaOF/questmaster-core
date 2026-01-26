package character

type CreateCharacterRequest struct {
	Name   string `json:"name"`
	System string `json:"system"`
	Hp     int    `json:"hp"`
}
