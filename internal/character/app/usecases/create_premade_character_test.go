package character

import (
	"errors"
	"testing"

	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakePremadeRepo struct {
	characterApp.CharacterRepository
	created   bool
	gotSystem rpgDomain.System
}

func (f *fakePremadeRepo) CreatePremade(name characterDomain.CharacterName, campaignID campaignDomain.CampaignID, system rpgDomain.System) (characterDomain.Character, error) {
	f.created, f.gotSystem = true, system
	return characterDomain.Character{Name: name, Slug: "harvey", CampaignID: &campaignID, System: system}, nil
}

type fakeFailingCampaignFinder struct{ err error }

func (f fakeFailingCampaignFinder) FindByID(campaignDomain.CampaignID) (campaignDomain.Campaign, error) {
	return campaignDomain.Campaign{}, f.err
}

func TestCreatePremadeCharacter(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())
	campaignIn := func(status campaignDomain.CampaignStatus) fakeCampaignFinder {
		return fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: 1, Dm: dm, Status: status, System: rpgDomain.CallOfCthulhu}}
	}
	cmd := func(user userDomain.UserID) characterApp.CreatePremadeCharacterCommand {
		return characterApp.CreatePremadeCharacterCommand{CampaignID: 1, Name: "Harvey", UserID: user}
	}

	t.Run("DM creates it with the campaign game system", func(t *testing.T) {
		repo := &fakePremadeRepo{}
		rm, err := NewCreatePremadeCharacter(repo, campaignIn(campaignDomain.StatusDraft)).Execute(cmd(dm))
		if err != nil || rm.Slug != "harvey" || repo.gotSystem != rpgDomain.CallOfCthulhu {
			t.Fatalf("expected a Call of Cthulhu premade, got %+v system=%v err=%v", rm, repo.gotSystem, err)
		}
	})

	t.Run("non-DM is forbidden", func(t *testing.T) {
		repo := &fakePremadeRepo{}
		_, err := NewCreatePremadeCharacter(repo, campaignIn(campaignDomain.StatusDraft)).Execute(cmd(userDomain.NewUserID(uuid.New())))
		if !errors.Is(err, campaignDomain.ErrNotDM) || repo.created {
			t.Fatalf("expected ErrNotDM without creating, got %v created=%v", err, repo.created)
		}
	})

	t.Run("archived campaign is rejected", func(t *testing.T) {
		repo := &fakePremadeRepo{}
		_, err := NewCreatePremadeCharacter(repo, campaignIn(campaignDomain.StatusArchived)).Execute(cmd(dm))
		if !errors.Is(err, campaignDomain.ErrCampaignArchived) || repo.created {
			t.Fatalf("expected ErrCampaignArchived without creating, got %v created=%v", err, repo.created)
		}
	})

	t.Run("unknown campaign", func(t *testing.T) {
		repo := &fakePremadeRepo{}
		_, err := NewCreatePremadeCharacter(repo, fakeFailingCampaignFinder{err: campaignUsecases.ErrCampaignNotFound}).Execute(cmd(dm))
		if !errors.Is(err, campaignUsecases.ErrCampaignNotFound) || repo.created {
			t.Fatalf("expected ErrCampaignNotFound without creating, got %v created=%v", err, repo.created)
		}
	})
}
