package user

import (
	"testing"

	userDomain "questmaster-core/internal/user/domain"

	"github.com/google/uuid"
)

func TestMapUserToUserResponseWithoutName(t *testing.T) {
	resp := MapUserToUserResponse(userDomain.User{
		Id:       userDomain.NewUserID(uuid.New()),
		Username: "player1",
	})

	if resp.Name != nil || resp.Surname != nil {
		t.Fatalf("expected nil name and surname, got %v %v", resp.Name, resp.Surname)
	}
}
