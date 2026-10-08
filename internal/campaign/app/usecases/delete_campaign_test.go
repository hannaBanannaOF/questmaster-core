package campaign

import (
	"errors"
	"testing"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeDeleteRepo struct {
	campaignApp.CampaignRepository
	campaign *campaignDomain.Campaign
	deleted  bool
}

func (f *fakeDeleteRepo) FindById(campaignDomain.CampaignID) (*campaignDomain.Campaign, error) {
	return f.campaign, nil
}

func (f *fakeDeleteRepo) DeleteById(campaignDomain.CampaignID) (bool, error) {
	f.deleted = true
	return true, nil
}

func TestDeleteCampaign(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())

	newRepo := func(status campaignDomain.CampaignStatus) *fakeDeleteRepo {
		return &fakeDeleteRepo{campaign: &campaignDomain.Campaign{Id: 1, Dm: dm, Status: status}}
	}

	for _, status := range []campaignDomain.CampaignStatus{campaignDomain.StatusDraft, campaignDomain.StatusArchived} {
		t.Run("DM deletes "+status.Value()+" campaign", func(t *testing.T) {
			repo := newRepo(status)
			if err := NewDeleteCampaign(repo).Execute(campaignApp.DeleteCampaignCommand{ID: 1, UserID: dm}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !repo.deleted {
				t.Fatalf("expected campaign to be deleted")
			}
		})
	}

	t.Run("non-DM is forbidden", func(t *testing.T) {
		repo := newRepo(campaignDomain.StatusDraft)
		other := userDomain.NewUserID(uuid.New())
		err := NewDeleteCampaign(repo).Execute(campaignApp.DeleteCampaignCommand{ID: 1, UserID: other})
		if !errors.Is(err, campaignDomain.ErrNotDM) || repo.deleted {
			t.Fatalf("expected ErrNotDM without deleting, got %v deleted=%v", err, repo.deleted)
		}
	})

	for _, status := range []campaignDomain.CampaignStatus{campaignDomain.StatusActive, campaignDomain.StatusPaused} {
		t.Run(status.Value()+" campaign is rejected", func(t *testing.T) {
			repo := newRepo(status)
			err := NewDeleteCampaign(repo).Execute(campaignApp.DeleteCampaignCommand{ID: 1, UserID: dm})
			if !errors.Is(err, campaignDomain.ErrNotDeletableStatus) || repo.deleted {
				t.Fatalf("expected ErrNotDeletableStatus without deleting, got %v deleted=%v", err, repo.deleted)
			}
		})
	}

	t.Run("campaign not found", func(t *testing.T) {
		err := NewDeleteCampaign(&fakeDeleteRepo{}).Execute(campaignApp.DeleteCampaignCommand{ID: 1, UserID: dm})
		if !errors.Is(err, ErrCampaignNotFound) {
			t.Fatalf("expected ErrCampaignNotFound, got %v", err)
		}
	})
}
