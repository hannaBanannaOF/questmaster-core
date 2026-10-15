package invite_test

import (
	"errors"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakePremadeFinder struct {
	premades []characterDomain.Character
}

func (f fakePremadeFinder) GetPremadeByCampaignID(campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	return f.premades, nil
}

func TestGetInviteDetail(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())
	hash := inviteDomain.NewHash(uuid.New())
	repo := fakeHashRepo{invite: &inviteDomain.Invite{CampaignId: 1, Hash: hash}}
	finder := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: 1, Name: "Lost Mine", Dm: dm}}

	details := func(r inviteApp.InviteRepository, user userDomain.UserID) (inviteApp.InviteDetailReadModel, error) {
		return inviteUsecases.NewGetInviteDetail(r, finder, fakePremadeFinder{}).Execute(inviteApp.GetinviteDetailsCommand{UserID: user, Hash: hash})
	}

	t.Run("player reads the invite", func(t *testing.T) {
		rm, err := details(repo, userDomain.NewUserID(uuid.New()))
		if err != nil || rm.IsDM || rm.CampaignName != "Lost Mine" {
			t.Fatalf("expected campaign details with IsDM false, got %+v err=%v", rm, err)
		}
	})

	t.Run("DM reads their own invite", func(t *testing.T) {
		rm, err := details(repo, dm)
		if err != nil || !rm.IsDM {
			t.Fatalf("expected IsDM true, got %+v err=%v", rm, err)
		}
	})

	t.Run("unknown invite", func(t *testing.T) {
		_, err := details(fakeHashRepo{}, dm)
		if !errors.Is(err, inviteUsecases.ErrInviteNotFound) {
			t.Fatalf("expected ErrInviteNotFound, got %v", err)
		}
	})

	t.Run("lists the unclaimed premade characters", func(t *testing.T) {
		premades := fakePremadeFinder{premades: []characterDomain.Character{{Slug: "harvey", Name: "Harvey"}}}
		rm, err := inviteUsecases.NewGetInviteDetail(repo, finder, premades).
			Execute(inviteApp.GetinviteDetailsCommand{UserID: dm, Hash: hash})
		if err != nil || len(rm.PremadeCharacters) != 1 {
			t.Fatalf("expected one premade character, got %+v err=%v", rm.PremadeCharacters, err)
		}
		if c := rm.PremadeCharacters[0]; c.Slug != "harvey" || c.Name != "Harvey" {
			t.Fatalf("unexpected premade character %+v", c)
		}
	})
}
