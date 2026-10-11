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

type fakeAccessRecorder struct {
	recorded []userDomain.UserID
	err      error
}

func (f *fakeAccessRecorder) RecordAccess(user userDomain.UserID, _ campaignDomain.CampaignID) error {
	f.recorded = append(f.recorded, user)
	return f.err
}

func newDetailsUC(dm, player userDomain.UserID, invites *fakeInviteFinder) *GetCampaignDetailsUseCase {
	return newDetailsUCWithRecorder(dm, player, invites, &fakeAccessRecorder{})
}

func newDetailsUCWithRecorder(dm, player userDomain.UserID, invites *fakeInviteFinder, recorder *fakeAccessRecorder) *GetCampaignDetailsUseCase {
	campaign := &campaignDomain.Campaign{Id: 1, Dm: dm, Status: campaignDomain.StatusActive}
	characters := []characterDomain.Character{{Id: 10, Name: "Hero", PlayerID: player}}
	return NewGetCampaignDetails(
		*NewGetCampaignFromID(fakeCampaignRepo{campaign: campaign}),
		fakeCharacterFinder{characters: characters},
		invites,
		recorder,
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

	t.Run("campaign not found", func(t *testing.T) {
		uc := NewGetCampaignDetails(
			*NewGetCampaignFromID(fakeCampaignRepo{}),
			fakeCharacterFinder{},
			&fakeInviteFinder{},
			&fakeAccessRecorder{},
		)
		_, err := uc.Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: dm})
		if !errors.Is(err, ErrCampaignNotFound) {
			t.Fatalf("expected ErrCampaignNotFound, got %v", err)
		}
	})

	t.Run("outsider is forbidden", func(t *testing.T) {
		uc := newDetailsUC(dm, player, &fakeInviteFinder{invite: invite})
		_, err := uc.Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: outsider})
		if !errors.Is(err, campaignDomain.ErrNotCampaignMember) {
			t.Fatalf("expected ErrNotCampaignMember, got %v", err)
		}
	})

	for name, user := range map[string]userDomain.UserID{"DM": dm, "player": player} {
		t.Run(name+" opening the campaign is recorded", func(t *testing.T) {
			recorder := &fakeAccessRecorder{}
			if _, err := newDetailsUCWithRecorder(dm, player, &fakeInviteFinder{}, recorder).
				Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: user}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(recorder.recorded) != 1 || recorder.recorded[0] != user {
				t.Fatalf("expected one access recorded for %s, got %v", name, recorder.recorded)
			}
		})
	}

	t.Run("outsider is not recorded", func(t *testing.T) {
		recorder := &fakeAccessRecorder{}
		_, _ = newDetailsUCWithRecorder(dm, player, &fakeInviteFinder{}, recorder).
			Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: outsider})
		if len(recorder.recorded) != 0 {
			t.Fatalf("expected nothing recorded, got %v", recorder.recorded)
		}
	})

	t.Run("failing recorder still returns the details", func(t *testing.T) {
		recorder := &fakeAccessRecorder{err: errors.New("database is down")}
		rm, err := newDetailsUCWithRecorder(dm, player, &fakeInviteFinder{}, recorder).
			Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: player})
		if err != nil || rm.Id != 1 {
			t.Fatalf("expected the details despite the recorder error, got %+v err=%v", rm, err)
		}
	})
}
