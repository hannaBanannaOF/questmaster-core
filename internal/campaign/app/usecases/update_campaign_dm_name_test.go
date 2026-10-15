package campaign

import (
	"errors"
	"testing"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeDmNameRepo struct {
	campaignApp.CampaignRepository
	campaign *campaignDomain.Campaign
	saved    *rpgDomain.TableName
}

func (f *fakeDmNameRepo) FindById(campaignDomain.CampaignID) (*campaignDomain.Campaign, error) {
	return f.campaign, nil
}

func (f *fakeDmNameRepo) UpdateDmName(_ campaignDomain.CampaignID, name rpgDomain.TableName) error {
	f.saved = &name
	return nil
}

func TestUpdateCampaignDmName(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())
	newRepo := func() *fakeDmNameRepo {
		return &fakeDmNameRepo{campaign: &campaignDomain.Campaign{Id: 1, Dm: dm}}
	}
	update := func(repo *fakeDmNameRepo, user userDomain.UserID) (rpgDomain.TableName, error) {
		return NewUpdateCampaignDmName(repo).Execute(campaignApp.UpdateCampaignDmNameCommand{CampaignID: 1, UserID: user, DmName: "Mestre Ana"})
	}

	t.Run("DM changes the name", func(t *testing.T) {
		repo := newRepo()
		got, err := update(repo, dm)
		if err != nil || got != "Mestre Ana" || repo.saved == nil || *repo.saved != "Mestre Ana" {
			t.Fatalf("expected the name saved, got %q err=%v saved=%v", got, err, repo.saved)
		}
	})

	t.Run("player is forbidden and nothing is saved", func(t *testing.T) {
		repo := newRepo()
		if _, err := update(repo, userDomain.NewUserID(uuid.New())); !errors.Is(err, campaignDomain.ErrNotDM) || repo.saved != nil {
			t.Fatalf("expected ErrNotDM without saving, got %v saved=%v", err, repo.saved)
		}
	})

	t.Run("unknown campaign", func(t *testing.T) {
		if _, err := update(&fakeDmNameRepo{}, dm); !errors.Is(err, ErrCampaignNotFound) {
			t.Fatalf("expected ErrCampaignNotFound, got %v", err)
		}
	})
}
