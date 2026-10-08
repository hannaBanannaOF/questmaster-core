package invite_test

import (
	"errors"
	"testing"

	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeInviteRepo struct {
	inviteApp.InviteRepository
	alreadyExists bool
	created       bool
}

func (f *fakeInviteRepo) Create(campaignID campaignDomain.CampaignID) (*inviteDomain.Invite, error) {
	if f.alreadyExists {
		return nil, nil
	}
	f.created = true
	return &inviteDomain.Invite{CampaignId: campaignID, Hash: inviteDomain.NewHash(uuid.New())}, nil
}

type fakeCampaignFinder struct {
	campaign campaignDomain.Campaign
	err      error
}

func (f fakeCampaignFinder) FindByID(campaignDomain.CampaignID) (campaignDomain.Campaign, error) {
	return f.campaign, f.err
}

func TestCreateInvite(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())
	finder := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: 1, Dm: dm}}

	t.Run("DM creates invite", func(t *testing.T) {
		repo := &fakeInviteRepo{}
		_, err := inviteUsecases.NewCreateInvite(repo, finder).Execute(inviteApp.CreateInviteCommand{CampaignID: 1, UserID: dm})
		if err != nil || !repo.created {
			t.Fatalf("expected invite to be created, err=%v", err)
		}
	})

	t.Run("non-DM is forbidden", func(t *testing.T) {
		repo := &fakeInviteRepo{}
		other := userDomain.NewUserID(uuid.New())
		_, err := inviteUsecases.NewCreateInvite(repo, finder).Execute(inviteApp.CreateInviteCommand{CampaignID: 1, UserID: other})
		if !errors.Is(err, campaignDomain.ErrNotDM) {
			t.Fatalf("expected ErrNotDM, got %v", err)
		}
		if repo.created {
			t.Fatalf("invite must not be created")
		}
	})

	t.Run("campaign not found", func(t *testing.T) {
		missing := fakeCampaignFinder{err: campaignUsecases.ErrCampaignNotFound}
		_, err := inviteUsecases.NewCreateInvite(&fakeInviteRepo{}, missing).Execute(inviteApp.CreateInviteCommand{CampaignID: 1, UserID: dm})
		if !errors.Is(err, campaignUsecases.ErrCampaignNotFound) {
			t.Fatalf("expected ErrCampaignNotFound, got %v", err)
		}
	})

	t.Run("archived campaign is rejected", func(t *testing.T) {
		repo := &fakeInviteRepo{}
		archived := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: 1, Dm: dm, Status: campaignDomain.StatusArchived}}
		_, err := inviteUsecases.NewCreateInvite(repo, archived).Execute(inviteApp.CreateInviteCommand{CampaignID: 1, UserID: dm})
		if !errors.Is(err, campaignDomain.ErrCampaignArchived) || repo.created {
			t.Fatalf("expected ErrCampaignArchived without creating, got %v created=%v", err, repo.created)
		}
	})

	t.Run("invite already exists", func(t *testing.T) {
		repo := &fakeInviteRepo{alreadyExists: true}
		_, err := inviteUsecases.NewCreateInvite(repo, finder).Execute(inviteApp.CreateInviteCommand{CampaignID: 1, UserID: dm})
		if !errors.Is(err, inviteUsecases.ErrInviteAlreadyExists) {
			t.Fatalf("expected ErrInviteAlreadyExists, got %v", err)
		}
	})
}
