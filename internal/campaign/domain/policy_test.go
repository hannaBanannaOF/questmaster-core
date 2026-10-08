package campaign

import (
	"errors"
	"testing"

	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

var allStatuses = []CampaignStatus{StatusDraft, StatusActive, StatusPaused, StatusArchived}

func TestCanInvite(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())

	for _, status := range allStatuses {
		t.Run(status.Value(), func(t *testing.T) {
			c := Campaign{Dm: dm, Status: status}
			err := c.CanInvite(dm)
			if status == StatusArchived {
				if !errors.Is(err, ErrCampaignArchived) {
					t.Fatalf("expected ErrCampaignArchived, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("expected DM to invite, got %v", err)
			}
		})
	}

	t.Run("non-DM", func(t *testing.T) {
		c := Campaign{Dm: dm, Status: StatusActive}
		if err := c.CanInvite(userDomain.NewUserID(uuid.New())); !errors.Is(err, ErrNotDM) {
			t.Fatalf("expected ErrNotDM, got %v", err)
		}
	})
}

func TestCanJoin(t *testing.T) {
	dm := userDomain.NewUserID(uuid.New())
	player := userDomain.NewUserID(uuid.New())

	for _, status := range allStatuses {
		t.Run(status.Value(), func(t *testing.T) {
			c := Campaign{Dm: dm, Status: status}
			err := c.CanJoin(player)
			if status == StatusArchived {
				if !errors.Is(err, ErrCampaignArchived) {
					t.Fatalf("expected ErrCampaignArchived, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("expected player to join, got %v", err)
			}
		})
	}

	t.Run("DM", func(t *testing.T) {
		c := Campaign{Dm: dm, Status: StatusActive}
		if err := c.CanJoin(dm); !errors.Is(err, ErrDMCannotJoin) {
			t.Fatalf("expected ErrDMCannotJoin, got %v", err)
		}
	})
}
