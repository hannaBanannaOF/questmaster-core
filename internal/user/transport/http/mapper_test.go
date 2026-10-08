package user

import (
	"testing"

	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

func TestMapUserToUserResponse(t *testing.T) {
	t.Run("with given name", func(t *testing.T) {
		surname := "Baggins"
		name, err := userDomain.NewName("Bilbo", &surname)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp := MapUserToUserResponse(userDomain.User{
			Id:       userDomain.NewUserID(uuid.New()),
			Username: "bilbo",
			Name:     &name,
		})
		if resp.Name == nil || *resp.Name != "Bilbo" || resp.Surname == nil || *resp.Surname != "Baggins" {
			t.Fatalf("expected name and surname from token, got %v %v", resp.Name, resp.Surname)
		}
	})

	t.Run("without given name", func(t *testing.T) {
		resp := MapUserToUserResponse(userDomain.User{
			Id:       userDomain.NewUserID(uuid.New()),
			Username: "player1",
		})
		if resp.Name != nil || resp.Surname != nil {
			t.Fatalf("expected nil name and surname, got %v %v", resp.Name, resp.Surname)
		}
	})
}
