package campaign

import (
	"testing"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeListRepo struct {
	campaignApp.CampaignRepository
	asDM, asPlayer []campaignDomain.Campaign
}

func (f fakeListRepo) GetByDmId(userDomain.UserID) ([]campaignDomain.Campaign, error) {
	return f.asDM, nil
}

func (f fakeListRepo) GetByPlayerId(userDomain.UserID) ([]campaignDomain.Campaign, error) {
	return f.asPlayer, nil
}

func TestGetCurrentUserCampaignsListsEachCampaignOnce(t *testing.T) {
	user := userDomain.NewUserID(uuid.New())
	mastered := campaignDomain.Campaign{Id: 1, Dm: user}
	played := campaignDomain.Campaign{Id: 2}

	// The same campaign coming from both queries must be listed once
	repo := fakeListRepo{
		asDM:     []campaignDomain.Campaign{mastered},
		asPlayer: []campaignDomain.Campaign{played, mastered},
	}

	list, err := NewGetCurrentUserMyCampaigns(repo).Execute(campaignApp.GetCurrentUserCampaignsCommand{UserID: user})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 || list[0].Id != mastered.Id || list[1].Id != played.Id {
		t.Fatalf("expected campaigns 1 and 2 once each, got %v", list)
	}
}
