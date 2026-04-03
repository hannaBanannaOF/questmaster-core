package character

import (
	characterUsecases "questmaster-core/internal/character/app/usecases"
	characterInfra "questmaster-core/internal/character/infra/pg"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CharacterModule struct {
	createCharacterUC                *characterUsecases.CreateCharacterUseCase
	deleteCharacterUC                *characterUsecases.DeleteCharacterUseCase
	getCurrentUserCharactersUC       *characterUsecases.GetCurrentUserCharactersUseCase
	getCampaignCharactersUC          *characterUsecases.GetCampaignCharactersUseCase
	getCharacterDetailUC             *characterUsecases.GetCharacterDetailUseCase
	resolveCharacterSlugUC           *characterUsecases.ResolveCharacterSlugUseCase
	updateHPUC                       *characterUsecases.UpdateHPUseCase
	getMyCharactersWithoutCampaignUC *characterUsecases.GetMyCharactersWithoutCampaignUseCase
	linkCharacterToCampaingUC        *characterUsecases.LinkCharacterToCampaignUseCase
}

func NewCharacterModule(db *pgxpool.Pool, campaignFinder characterUsecases.CharacterCampaingFinder) *CharacterModule {
	r := characterInfra.NewCharacterRepositoryPG(db)
	return &CharacterModule{
		createCharacterUC:                characterUsecases.NewCreateCharacter(r),
		deleteCharacterUC:                characterUsecases.NewDeleteCharacter(r),
		getCurrentUserCharactersUC:       characterUsecases.NewGetCurrrentUserCharacters(r),
		getCampaignCharactersUC:          characterUsecases.NewGetCampaignCharacters(r),
		getCharacterDetailUC:             characterUsecases.NewGetCharacterDetail(r),
		resolveCharacterSlugUC:           characterUsecases.NewResolveCharacterSlug(r),
		updateHPUC:                       characterUsecases.NewUpdateHP(r, campaignFinder),
		getMyCharactersWithoutCampaignUC: characterUsecases.NewMyGetCharactersWithoutCampaign(r),
		linkCharacterToCampaingUC:        characterUsecases.NewLinkCharacterToCampaign(r),
	}
}

func (m *CharacterModule) CreateCharacterUC() *characterUsecases.CreateCharacterUseCase {
	return m.createCharacterUC
}

func (m *CharacterModule) DeleteCharacterUC() *characterUsecases.DeleteCharacterUseCase {
	return m.deleteCharacterUC
}

func (m *CharacterModule) GetCurrentUserCharactersUC() *characterUsecases.GetCurrentUserCharactersUseCase {
	return m.getCurrentUserCharactersUC
}

func (m *CharacterModule) GetCampaignCharactersUC() *characterUsecases.GetCampaignCharactersUseCase {
	return m.getCampaignCharactersUC
}

func (m *CharacterModule) GetCharacterDetailUC() *characterUsecases.GetCharacterDetailUseCase {
	return m.getCharacterDetailUC
}

func (m *CharacterModule) ResolveCharacterSlugUC() *characterUsecases.ResolveCharacterSlugUseCase {
	return m.resolveCharacterSlugUC
}

func (m *CharacterModule) UpdateHPUC() *characterUsecases.UpdateHPUseCase {
	return m.updateHPUC
}

func (m *CharacterModule) GetMyCharactersWithoutCampaignUC() *characterUsecases.GetMyCharactersWithoutCampaignUseCase {
	return m.getMyCharactersWithoutCampaignUC
}

func (m *CharacterModule) LinkCharacterToCampaignUC() *characterUsecases.LinkCharacterToCampaignUseCase {
	return m.linkCharacterToCampaingUC
}
