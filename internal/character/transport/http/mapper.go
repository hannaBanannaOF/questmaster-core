package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	rpgTransport "questmaster-core/internal/rpg/transport/http"
	"questmaster-core/internal/shared/pagination"
	userDomain "questmaster-core/internal/user/domain"
)

func MapListReadModelToResponse(c characterDomain.Character) CharacterListResponse {
	var campaign *CharacterCampaignResponse
	if c.Campaign != nil {
		campaign = &CharacterCampaignResponse{Slug: c.Campaign.Slug, Name: c.Campaign.Name}
	}
	return CharacterListResponse{
		Campaign:       campaign,
		Slug:           c.Slug.Value(),
		Name:           c.Name.Value(),
		System:         c.System.Value(),
		LastAccessedAt: c.LastAccessedAt,
	}
}

func MapDetailReadModelToResponse(c characterDomain.Character, userID userDomain.UserID) CharacterDetailResponse {
	return CharacterDetailResponse{
		Id:       c.Id.Value(),
		Name:     c.Name.Value(),
		System:   c.System.Value(),
		Slug:     c.Slug.Value(),
		IsPlayer: c.IsPlayer(userID),
	}
}

func MapCreateCharacterReadModelToResponse(rm characterApp.CreateCharacterReadModel) rpgTransport.RpgSlugResponse {
	return rpgTransport.RpgSlugResponse{
		Slug: rm.Slug,
	}
}

func MapResolveSlugReadModelToResponse(rm characterApp.CharacterResolveSlugReadModel) rpgTransport.RpgIdResponse {
	return rpgTransport.RpgIdResponse{
		ID: rm.ID,
	}
}

func MapListPageToResponse(page pagination.Result[characterDomain.Character]) CharacterListPageResponse {
	items := make([]CharacterListResponse, 0, len(page.Items))
	for _, c := range page.Items {
		items = append(items, MapListReadModelToResponse(c))
	}

	return CharacterListPageResponse{
		Items: items,
		Total: page.Total,
	}
}
