package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgTransport "questmaster-core/internal/rpg/transport/http"
	"questmaster-core/internal/shared/pagination"
	user "questmaster-core/internal/user/domain"
)

func MapListReadModelToResponse(c campaignDomain.Campaign, userID user.UserID) CampaignListResponse {
	return CampaignListResponse{
		Slug:        c.Slug.Value(),
		Name:        c.Name.Value(),
		IsDM:        c.IsDM(userID),
		Status:      c.Status.Value(),
		System:      c.System.Value(),
		PlayerCount: c.PlayerCount.Value(),
	}
}

func MapStatusToResponse(status string) CampaignStatusResponse {
	return CampaignStatusResponse{
		Status: status,
	}
}

func MapDetailReadModelToResponse(rm campaignApp.CampaignDetailsReadModel) CampaignDetailResponse {

	characters := make([]CampaignDetailResponseCharacterItem, 0, len(rm.Characters))

	for _, c := range rm.Characters {
		characters = append(characters, CampaignDetailResponseCharacterItem{
			Id:        c.Id,
			Name:      c.Name,
			CurrentHP: c.CurrentHP,
		})
	}

	return CampaignDetailResponse{
		Id:         rm.Id,
		Name:       rm.Name,
		Status:     rm.Status,
		System:     rm.System,
		Slug:       rm.Slug,
		Overview:   rm.Overview,
		IsDM:       rm.IsDM,
		Characters: characters,
		InviteHash: rm.InviteHash,
	}
}

func MapResolveSlugReadModelToResponse(rm campaignApp.ResolveCampaignSlugReadModel) rpgTransport.RpgIdResponse {
	return rpgTransport.RpgIdResponse{
		ID: rm.ID,
	}
}

func MapCreateCampaignReadModelToResponse(rm campaignApp.CreateCampaignReadModel) rpgTransport.RpgSlugResponse {
	return rpgTransport.RpgSlugResponse{
		Slug: rm.Slug,
	}
}

func MapListPageToResponse(page pagination.Result[campaignDomain.Campaign], userID user.UserID) CampaignListPageResponse {
	items := make([]CampaignListResponse, 0, len(page.Items))
	for _, c := range page.Items {
		items = append(items, MapListReadModelToResponse(c, userID))
	}

	return CampaignListPageResponse{
		Items: items,
		Total: page.Total,
	}
}

func MapStatusCountsToResponse(counts campaignApp.CampaignStatusCountsReadModel) CampaignStatusCountsResponse {
	return CampaignStatusCountsResponse{
		Draft:    counts[campaignDomain.StatusDraft.Value()],
		Active:   counts[campaignDomain.StatusActive.Value()],
		Paused:   counts[campaignDomain.StatusPaused.Value()],
		Archived: counts[campaignDomain.StatusArchived.Value()],
	}
}
