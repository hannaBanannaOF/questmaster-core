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

type UpdateHPReadModel struct {
	CurrentHP int
}

type CharacterResolveSlugReadModel struct {
	ID int
}

type CreateCharacterReadModel struct {
	Slug string
}
