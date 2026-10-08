package invite

import (
	"errors"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeInviteRepo struct {
	inviteApp.InviteRepository
	created bool
}

func (f *fakeInviteRepo) Create(campaignID campaignDomain.CampaignID) (*inviteDomain.Invite, error) {
	f.created = true
	return &inviteDomain.Invite{CampaignId: campaignID, Hash: inviteDomain.NewHash(uuid.New())}, nil
}

type fakeCampaignFinder struct {
	campaign campaignDomain.Campaign
}

func (f fakeCampaignFinder) FindByID(campaignDomain.CampaignID) (campaignDomain.Campaign, error) {
	return f.campaign, nil
}

func TestCreateInvite(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())
	finder := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: 1, Dm: dm}}

	t.Run("DM creates invite", func(t *testing.T) {
		repo := &fakeInviteRepo{}
		_, err := NewCreateInvite(repo, finder).Execute(inviteApp.CreateInviteCommand{CampaignID: 1, UserID: dm})
		if err != nil || !repo.created {
			t.Fatalf("expected invite to be created, err=%v", err)
		}
	})

	t.Run("non-DM is forbidden", func(t *testing.T) {
		repo := &fakeInviteRepo{}
		other := userDomain.NewUserID(uuid.New())
		_, err := NewCreateInvite(repo, finder).Execute(inviteApp.CreateInviteCommand{CampaignID: 1, UserID: other})
		if !errors.Is(err, campaignDomain.ErrNotDM) {
			t.Fatalf("expected ErrNotDM, got %v", err)
		}
		if repo.created {
			t.Fatalf("invite must not be created")
		}
	})
}
