package character

type CharacterResolveSlugReadModel struct {
	ID int
}

type EditCharacterReadModel struct {
	Slug       string
	Name       string
	PlayerName *string
}

type CreateCharacterReadModel struct {
	Slug string
}
