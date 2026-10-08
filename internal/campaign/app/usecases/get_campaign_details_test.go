package campaign

import (
	"errors"
	"testing"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeCampaignRepo struct {
	campaignApp.CampaignRepository
	campaign *campaignDomain.Campaign
}

func (f fakeCampaignRepo) FindById(campaignDomain.CampaignID) (*campaignDomain.Campaign, error) {
	return f.campaign, nil
}

type fakeCharacterFinder struct {
	characters []characterDomain.Character
}

func (f fakeCharacterFinder) GetByCampaignID(campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	return f.characters, nil
}

type fakeInviteFinder struct {
	invite *inviteDomain.Invite
	called bool
}

func (f *fakeInviteFinder) GetByCampaignID(campaignDomain.CampaignID) (*inviteDomain.Invite, error) {
	f.called = true
	return f.invite, nil
}

func newDetailsUC(dm, player userDomain.UserID, invites *fakeInviteFinder) *GetCampaignDetailsUseCase {
	campaign := &campaignDomain.Campaign{Id: 1, Dm: dm, Status: campaignDomain.StatusActive}
	characters := []characterDomain.Character{{Id: 10, Name: "Hero", PlayerID: player}}
	return NewGetCampaignDetails(
		*NewGetCampaignFromID(fakeCampaignRepo{campaign: campaign}),
		fakeCharacterFinder{characters: characters},
		invites,
	)
}

func TestGetCampaignDetails(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())
	player := userDomain.NewUserID(uuid.New())
	outsider := userDomain.NewUserID(uuid.New())
	invite := &inviteDomain.Invite{Hash: inviteDomain.NewHash(uuid.New())}

	t.Run("DM sees invite hash", func(t *testing.T) {
		uc := newDetailsUC(dm, player, &fakeInviteFinder{invite: invite})
		rm, err := uc.Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: dm})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rm.InviteHash == nil || *rm.InviteHash != invite.Hash.Value() {
			t.Fatalf("expected DM to get invite hash, got %v", rm.InviteHash)
		}
	})

	t.Run("player sees details without invite hash", func(t *testing.T) {
		invites := &fakeInviteFinder{invite: invite}
		uc := newDetailsUC(dm, player, invites)
		rm, err := uc.Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: player})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rm.InviteHash != nil || invites.called {
			t.Fatalf("player must not get invite hash")
		}
	})

	t.Run("outsider is forbidden", func(t *testing.T) {
		uc := newDetailsUC(dm, player, &fakeInviteFinder{invite: invite})
		_, err := uc.Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: outsider})
		if !errors.Is(err, campaignDomain.ErrNotCampaignMember) {
			t.Fatalf("expected ErrNotCampaignMember, got %v", err)
		}
	})
}
