package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	rpgTransport "questmaster-core/internal/rpg/transport/http"
)

func MapListReadModelToResponse(rm campaignApp.CampaignListReadModel) CampaignListResponse {
	return CampaignListResponse{
		Slug:   rm.Slug,
		Name:   rm.Name,
		IsDM:   rm.IsDM,
		Status: rm.Status,
		System: rm.System,
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
			Id:   c.Id,
			Name: c.Name,
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
	}
}

func MapCreateRequestToCreateCommand(req CreateCampaignRequest, userID rpgDomain.UserID) (campaignApp.CreateCampaignCommand, error) {
	name, err := campaignDomain.NewCampaignName(req.Name)
	if err != nil {
		return campaignApp.CreateCampaignCommand{}, err
	}

	var overview *campaignDomain.CampaignOverview

	if req.Overview != nil {
		o := campaignDomain.NewCampaignOverview(*req.Overview)
		overview = &o
	}

	system, err := rpgDomain.NewSystem(req.System)
	if err != nil {
		return campaignApp.CreateCampaignCommand{}, err
	}

	return campaignApp.CreateCampaignCommand{
		Name:     name,
		Overview: overview,
		DmID:     userID,
		System:   system,
	}, nil
}

func MapGetOrCreateCampaignInviteReadModelToResponse(rm campaignApp.GetOrCreateInviteReadModel) CreateCampaignInviteResponse {
	return CreateCampaignInviteResponse{
		Hash: rm.InviteHash,
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
