package campaign

import (
	"testing"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeCountRepo struct {
	campaignApp.CampaignRepository
	counts map[campaignDomain.CampaignStatus]int
}

func (f fakeCountRepo) CountByStatusForUser(userDomain.UserID, campaignDomain.CampaignListFilters) (map[campaignDomain.CampaignStatus]int, error) {
	return f.counts, nil
}

func TestGetCampaignStatusCounts(t *testing.T) {
	user := userDomain.NewUserID(uuid.New())

	t.Run("user without campaigns gets every status with 0", func(t *testing.T) {
		got, err := NewGetCampaignStatusCounts(fakeCountRepo{counts: map[campaignDomain.CampaignStatus]int{}}).
			Execute(campaignApp.GetCampaignStatusCountsCommand{UserID: user})
		if err != nil || len(got) != 4 || got["DRAFT"] != 0 || got["ACTIVE"] != 0 || got["PAUSED"] != 0 || got["ARCHIVED"] != 0 {
			t.Fatalf("expected the 4 statuses with 0, got %v err=%v", got, err)
		}
	})

	t.Run("missing statuses are filled with 0", func(t *testing.T) {
		repo := fakeCountRepo{counts: map[campaignDomain.CampaignStatus]int{
			campaignDomain.StatusActive:   2,
			campaignDomain.StatusArchived: 1,
		}}
		got, err := NewGetCampaignStatusCounts(repo).Execute(campaignApp.GetCampaignStatusCountsCommand{UserID: user})
		if err != nil || got["ACTIVE"] != 2 || got["ARCHIVED"] != 1 || got["DRAFT"] != 0 || got["PAUSED"] != 0 {
			t.Fatalf("unexpected counts %v err=%v", got, err)
		}
	})
}
