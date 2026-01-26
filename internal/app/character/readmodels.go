package character

type CharacterListReadModel struct {
	Slug   string
	Name   string
	System string
}

type CharacterDetailReadModel struct {
	Name      string
	MaxHp     int
	CurrentHp int
}
