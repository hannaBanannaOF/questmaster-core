package campaign

import (
	"errors"
	"testing"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
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

type fakePremadeFinder struct {
	premades []characterDomain.Character
}

func (f fakePremadeFinder) GetPremadeByCampaignID(campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	return f.premades, nil
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
	characters := []characterDomain.Character{{Id: 10, Name: "Hero", PlayerID: &player}}
	return NewGetCampaignDetails(
		*NewGetCampaignFromID(fakeCampaignRepo{campaign: campaign}),
		fakeCharacterFinder{characters: characters},
		invites,
		fakePremadeFinder{},
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
			fakePremadeFinder{},
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

	t.Run("party cards", func(t *testing.T) {
		other := userDomain.NewUserID(uuid.New())
		ana := rpgDomain.TableName("Ana")
		characters := []characterDomain.Character{
			{Id: 10, Slug: "hero", Name: "Hero", PlayerID: &player, PlayerName: &ana},
			{Id: 11, Slug: "sidekick", Name: "Sidekick", PlayerID: &other},
		}
		uc := NewGetCampaignDetails(
			*NewGetCampaignFromID(fakeCampaignRepo{campaign: &campaignDomain.Campaign{Id: 1, Dm: dm, DmName: "Mestre"}}),
			fakeCharacterFinder{characters: characters},
			&fakeInviteFinder{},
			fakePremadeFinder{},
			&fakeAccessRecorder{},
		)

		rm, err := uc.Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: player})
		if err != nil || len(rm.Characters) != 2 {
			t.Fatalf("unexpected details %+v err=%v", rm, err)
		}
		mine, others := rm.Characters[0], rm.Characters[1]
		if rm.DmName != "Mestre" || mine.PlayerName == nil || *mine.PlayerName != "Ana" || others.PlayerName != nil {
			t.Fatalf("expected the table names, got dm=%q %v %v", rm.DmName, mine.PlayerName, others.PlayerName)
		}
		if !mine.IsMine || mine.Slug != "hero" {
			t.Fatalf("expected own character marked with slug, got %+v", mine)
		}
		if others.IsMine || others.Slug != "sidekick" {
			t.Fatalf("expected other character not mine, got %+v", others)
		}

		rm, err = uc.Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: dm})
		if err != nil || rm.Characters[0].IsMine || rm.Characters[1].IsMine || rm.Characters[0].Slug != "hero" {
			t.Fatalf("expected the DM to see the party with is_mine false, got %+v err=%v", rm.Characters, err)
		}
	})

	t.Run("only the DM gets the premade characters", func(t *testing.T) {
		premades := fakePremadeFinder{premades: []characterDomain.Character{{Id: 20, Slug: "harvey", Name: "Harvey"}}}
		uc := NewGetCampaignDetails(
			*NewGetCampaignFromID(fakeCampaignRepo{campaign: &campaignDomain.Campaign{Id: 1, Dm: dm}}),
			fakeCharacterFinder{characters: []characterDomain.Character{{Id: 10, Name: "Hero", PlayerID: &player}}},
			&fakeInviteFinder{},
			premades,
			&fakeAccessRecorder{},
		)

		rm, err := uc.Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: dm})
		if err != nil || len(rm.PremadeCharacters) != 1 || rm.PremadeCharacters[0].Slug != "harvey" || len(rm.Characters) != 1 {
			t.Fatalf("expected the DM to get the party and the premade apart, got %+v err=%v", rm, err)
		}

		rm, err = uc.Execute(campaignApp.GetCampaignDetailsCommand{ID: 1, UserID: player})
		if err != nil || len(rm.PremadeCharacters) != 0 {
			t.Fatalf("expected no premade characters for a player, got %+v err=%v", rm.PremadeCharacters, err)
		}
	})
}
