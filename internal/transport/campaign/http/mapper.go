package campaign

import app "questmaster-core/internal/app/campaign"

func MapListReadModelToResponse(rm app.CampaignListReadModel) CampaignListResponse {
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
