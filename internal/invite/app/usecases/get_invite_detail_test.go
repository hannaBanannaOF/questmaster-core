package invite_test

import (
	"errors"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

func TestGetInviteDetail(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())
	hash := inviteDomain.NewHash(uuid.New())
	repo := fakeHashRepo{invite: &inviteDomain.Invite{CampaignId: 1, Hash: hash}}
	finder := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: 1, Name: "Lost Mine", Dm: dm}}

	details := func(r inviteApp.InviteRepository, user userDomain.UserID) (inviteApp.InviteDetailReadModel, error) {
		return inviteUsecases.NewGetInviteDetail(r, finder).Execute(inviteApp.GetinviteDetailsCommand{UserID: user, Hash: hash})
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
}
