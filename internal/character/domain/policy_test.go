package character

import (
	"errors"
	"testing"

	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

type fakeCampaign struct{ dm userDomain.UserID }

func (f fakeCampaign) IsDM(user userDomain.UserID) bool { return f.dm == user }

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
