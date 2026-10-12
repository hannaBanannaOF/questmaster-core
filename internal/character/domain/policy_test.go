package character

import (
	"errors"
	"testing"

	campaignDomain "questmaster-core/internal/campaign/domain"
	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeCampaign struct {
	dm       userDomain.UserID
	archived bool
}

func (f fakeCampaign) IsDM(user userDomain.UserID) bool { return f.dm == user }
func (f fakeCampaign) IsArchived() bool                 { return f.archived }

func TestCanView(t *testing.T) {
	player, dm, member, outsider := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New()),
		userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	c := Character{PlayerID: &player}
	campaign := fakeCampaign{dm: dm}

	allowed := []struct {
		name     string
		user     userDomain.UserID
		campaign CampaignAccess
		member   bool
	}{
		{"player", player, campaign, false},
		{"player without campaign", player, nil, false},
		{"campaign DM", dm, campaign, false},
		{"campaign member", member, campaign, true},
	}
	for _, tc := range allowed {
		if err := c.CanView(tc.user, tc.campaign, tc.member); err != nil {
			t.Fatalf("%s: expected allowed, got %v", tc.name, err)
		}
	}

	if err := c.CanView(outsider, campaign, false); !errors.Is(err, ErrCannotViewCharacter) {
		t.Fatalf("non-member: expected ErrCannotViewCharacter, got %v", err)
	}
	if err := c.CanView(outsider, nil, false); !errors.Is(err, ErrCannotViewCharacter) {
		t.Fatalf("no campaign: expected ErrCannotViewCharacter, got %v", err)
	}
}

func TestCanDelete(t *testing.T) {
	player, dm, other := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	campaign := fakeCampaign{dm: dm}
	owned := Character{PlayerID: &player}
	premade := Character{}

	if err := owned.CanDelete(player, campaign); err != nil {
		t.Fatalf("player: expected allowed, got %v", err)
	}
	if err := premade.CanDelete(dm, campaign); err != nil {
		t.Fatalf("DM on unclaimed premade: expected allowed, got %v", err)
	}
	for name, tc := range map[string]struct {
		c    Character
		user userDomain.UserID
	}{
		"DM on a claimed character":  {owned, dm},
		"other user on a premade":    {premade, other},
		"other user on an owned one": {owned, other},
	} {
		if err := tc.c.CanDelete(tc.user, campaign); !errors.Is(err, ErrNotPlayer) {
			t.Fatalf("%s: expected ErrNotPlayer, got %v", name, err)
		}
	}
	if err := premade.CanDelete(dm, nil); !errors.Is(err, ErrNotPlayer) {
		t.Fatalf("premade without campaign: expected ErrNotPlayer, got %v", err)
	}
}

func TestCanLeaveCampaign(t *testing.T) {
	player, dm, member := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	active := fakeCampaign{dm: dm}
	owned := Character{PlayerID: &player}

	for name, user := range map[string]userDomain.UserID{"player": player, "campaign DM": dm} {
		if err := owned.CanLeaveCampaign(user, active); err != nil {
			t.Fatalf("%s: expected allowed, got %v", name, err)
		}
	}
	if err := owned.CanLeaveCampaign(member, active); !errors.Is(err, ErrNotAllowed) {
		t.Fatalf("other member: expected ErrNotAllowed, got %v", err)
	}
	premade := Character{}
	if err := premade.CanLeaveCampaign(dm, active); !errors.Is(err, ErrPremadeCannotLeave) {
		t.Fatalf("premade: expected ErrPremadeCannotLeave, got %v", err)
	}
	if err := owned.CanLeaveCampaign(player, fakeCampaign{dm: dm, archived: true}); !errors.Is(err, campaignDomain.ErrCampaignArchived) {
		t.Fatalf("archived: expected ErrCampaignArchived, got %v", err)
	}
}
