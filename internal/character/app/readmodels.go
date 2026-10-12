package character

type UpdateHPReadModel struct {
	CurrentHP int
}

type CharacterResolveSlugReadModel struct {
	ID int
}

type EditCharacterReadModel struct {
	Slug      string
	Name      string
	CurrentHP *int
	MaxHP     *int
}

type CreateCharacterReadModel struct {
	Slug string
}
