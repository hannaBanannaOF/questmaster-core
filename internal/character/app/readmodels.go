package character

type UpdateHPReadModel struct {
	CurrentHP int
}

type CharacterResolveSlugReadModel struct {
	ID int
}

type CreateCharacterReadModel struct {
	Slug string
}
