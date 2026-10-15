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

// fakePlayerNames has one player with a character in the campaign
type fakePlayerNames struct {
	player userDomain.UserID
	saved  *rpgDomain.TableName
}

func (f *fakePlayerNames) UpdatePlayerName(user userDomain.UserID, _ campaignDomain.CampaignID, name rpgDomain.TableName) (bool, error) {
	if user != f.player {
		return false, nil
	}
	f.saved = &name
	return true, nil
}

func TestUpdateTableName(t *testing.T) {
	dm, player := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	update := func(repo *fakeDmNameRepo, user userDomain.UserID) (rpgDomain.TableName, *fakePlayerNames, error) {
		players := &fakePlayerNames{player: player}
		got, err := NewUpdateTableName(repo, players).Execute(campaignApp.UpdateTableNameCommand{CampaignID: 1, UserID: user, Name: "Ana"})
		return got, players, err
	}
	newRepo := func() *fakeDmNameRepo {
		return &fakeDmNameRepo{campaign: &campaignDomain.Campaign{Id: 1, Dm: dm}}
	}

	t.Run("DM changes the DM name", func(t *testing.T) {
		repo := newRepo()
		got, players, err := update(repo, dm)
		if err != nil || got != "Ana" || repo.saved == nil || *repo.saved != "Ana" || players.saved != nil {
			t.Fatalf("expected only dm_name saved, got %q err=%v dm=%v player=%v", got, err, repo.saved, players.saved)
		}
	})

	t.Run("player changes their player name", func(t *testing.T) {
		repo := newRepo()
		got, players, err := update(repo, player)
		if err != nil || got != "Ana" || players.saved == nil || *players.saved != "Ana" || repo.saved != nil {
			t.Fatalf("expected only player_name saved, got %q err=%v dm=%v player=%v", got, err, repo.saved, players.saved)
		}
	})

	t.Run("outsider is forbidden", func(t *testing.T) {
		repo := newRepo()
		_, players, err := update(repo, userDomain.NewUserID(uuid.New()))
		if !errors.Is(err, campaignDomain.ErrNotCampaignMember) || repo.saved != nil || players.saved != nil {
			t.Fatalf("expected ErrNotCampaignMember without saving, got %v", err)
		}
	})

	t.Run("unknown campaign", func(t *testing.T) {
		if _, _, err := update(&fakeDmNameRepo{}, dm); !errors.Is(err, ErrCampaignNotFound) {
			t.Fatalf("expected ErrCampaignNotFound, got %v", err)
		}
	})
}
