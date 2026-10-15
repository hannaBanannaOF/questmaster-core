package character

type CharacterResolveSlugReadModel struct {
	ID int
}

type EditCharacterReadModel struct {
	Slug string
	Name string
}

type CreateCharacterReadModel struct {
	Slug string
}
