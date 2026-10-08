package invite_test

import (
	"errors"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
	inviteDomain "questmaster-core/internal/invite/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeHashRepo struct {
	inviteApp.InviteRepository
	invite *inviteDomain.Invite
}

func (f fakeHashRepo) FindByHash(inviteDomain.InviteHash) (*inviteDomain.Invite, error) {
	return f.invite, nil
}

type fakeLinker struct {
	linked bool
}

func (f *fakeLinker) LinkToCampaign(campaignDomain.CampaignID, rpgDomain.Slug, userDomain.UserID) (characterDomain.Character, error) {
	f.linked = true
	return characterDomain.Character{}, nil
}

func TestAcceptInvite(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())
	player := userDomain.NewUserID(uuid.New())
	hash := inviteDomain.NewHash(uuid.New())
	repo := fakeHashRepo{invite: &inviteDomain.Invite{CampaignId: 1, Hash: hash}}

	accept := func(status campaignDomain.CampaignStatus, user userDomain.UserID, r inviteApp.InviteRepository) (*fakeLinker, error) {
		linker := &fakeLinker{}
		finder := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: 1, Dm: dm, Status: status}}
		err := inviteUsecases.NewAcceptInvite(r, finder, linker).Execute(inviteApp.AcceptInviteCommand{
			Hash: hash, CharacterSlug: "hero", UserID: user,
		})
		return linker, err
	}

	for _, status := range []campaignDomain.CampaignStatus{campaignDomain.StatusDraft, campaignDomain.StatusActive, campaignDomain.StatusPaused} {
		t.Run("player joins "+status.Value()+" campaign", func(t *testing.T) {
			linker, err := accept(status, player, repo)
			if err != nil || !linker.linked {
				t.Fatalf("expected character to be linked, err=%v", err)
			}
		})
	}

	t.Run("archived campaign is rejected", func(t *testing.T) {
		linker, err := accept(campaignDomain.StatusArchived, player, repo)
		if !errors.Is(err, campaignDomain.ErrCampaignArchived) || linker.linked {
			t.Fatalf("expected ErrCampaignArchived without linking, got %v linked=%v", err, linker.linked)
		}
	})

	t.Run("DM is forbidden", func(t *testing.T) {
		linker, err := accept(campaignDomain.StatusActive, dm, repo)
		if !errors.Is(err, campaignDomain.ErrDMCannotJoin) || linker.linked {
			t.Fatalf("expected ErrDMCannotJoin without linking, got %v linked=%v", err, linker.linked)
		}
	})

	t.Run("unknown invite", func(t *testing.T) {
		linker, err := accept(campaignDomain.StatusActive, player, fakeHashRepo{})
		if !errors.Is(err, inviteUsecases.ErrInviteNotFound) || linker.linked {
			t.Fatalf("expected ErrInviteNotFound without linking, got %v linked=%v", err, linker.linked)
		}
	})
}
