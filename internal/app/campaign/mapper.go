package campaign

import (
	domain "questmaster-core/internal/domain/campaign"

	"github.com/google/uuid"
)

func MapDomainToListReadModel(domain domain.Campaign, userId uuid.UUID) CampaignListReadModel {
	return CampaignListReadModel{
		Slug:   domain.Slug.String(),
		Name:   domain.Name.String(),
		System: string(domain.System),
		IsDM:   domain.IsDM(userId),
		Status: string(domain.Status),
	}
}
