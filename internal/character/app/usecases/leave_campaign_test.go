package character

import (
	"errors"
	"testing"

	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeLeaveRepo struct {
	characterApp.CharacterRepository
	character *characterDomain.Character
	unlinkOK  bool
	unlinked  bool
}

func (f *fakeLeaveRepo) FindByID(characterDomain.CharacterID) (*characterDomain.Character, error) {
	return f.character, nil
}

func (f *fakeLeaveRepo) UnlinkFromCampaign(characterDomain.CharacterID, campaignDomain.CampaignID) (bool, error) {
	f.unlinked = f.unlinkOK
	return f.unlinkOK, nil
}

func TestLeaveCampaign(t *testing.T) {
	player, dm, member := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	campaignID, otherCampaign := campaignDomain.NewCampaignID(1), campaignDomain.NewCampaignID(2)
	active := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: campaignID, Dm: dm, Status: campaignDomain.StatusActive}}
	archived := fakeCampaignFinder{campaign: campaignDomain.Campaign{Id: campaignID, Dm: dm, Status: campaignDomain.StatusArchived}}
	owned := func(in campaignDomain.CampaignID) *characterDomain.Character {
		return &characterDomain.Character{Id: 1, PlayerID: &player, CampaignID: &in}
	}

	leave := func(c *characterDomain.Character, finder CharacterCampaingFinder, user userDomain.UserID) (*fakeLeaveRepo, error) {
		repo := &fakeLeaveRepo{character: c, unlinkOK: true}
		err := NewLeaveCampaign(repo, finder).Execute(characterApp.LeaveCampaignCommand{CampaignID: campaignID, CharacterID: 1, UserID: user})
		return repo, err
	}

	for name, user := range map[string]userDomain.UserID{"player leaves": player, "DM removes": dm} {
		t.Run(name, func(t *testing.T) {
			repo, err := leave(owned(campaignID), active, user)
			if err != nil || !repo.unlinked {
				t.Fatalf("expected the character unlinked, got err=%v unlinked=%v", err, repo.unlinked)
			}
		})
	}

	rejected := []struct {
		name      string
		character *characterDomain.Character
		finder    CharacterCampaingFinder
		user      userDomain.UserID
		want      error
	}{
		{"other member is forbidden", owned(campaignID), active, member, characterDomain.ErrNotAllowed},
		{"unknown character", nil, active, player, ErrCharacterNotFound},
		{"unknown campaign", owned(campaignID), fakeFailingCampaignFinder{err: campaignUsecases.ErrCampaignNotFound}, player, campaignUsecases.ErrCampaignNotFound},
		{"character of another campaign", owned(otherCampaign), active, dm, ErrNotInCampaign},
		{"unclaimed premade character", &characterDomain.Character{Id: 1, CampaignID: &campaignID}, active, dm, characterDomain.ErrPremadeCannotLeave},
		{"archived campaign", owned(campaignID), archived, player, campaignDomain.ErrCampaignArchived},
	}
	for _, tc := range rejected {
		t.Run(tc.name, func(t *testing.T) {
			repo, err := leave(tc.character, tc.finder, tc.user)
			if !errors.Is(err, tc.want) || repo.unlinked {
				t.Fatalf("expected %v without unlinking, got %v unlinked=%v", tc.want, err, repo.unlinked)
			}
		})
	}

	t.Run("removed meanwhile", func(t *testing.T) {
		repo := &fakeLeaveRepo{character: owned(campaignID), unlinkOK: false}
		err := NewLeaveCampaign(repo, active).Execute(characterApp.LeaveCampaignCommand{CampaignID: campaignID, CharacterID: 1, UserID: player})
		if !errors.Is(err, ErrNotInCampaign) {
			t.Fatalf("expected ErrNotInCampaign, got %v", err)
		}
	})
}
