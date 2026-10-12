package character

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	campaignDomain "questmaster-core/internal/campaign/domain"
	campaignTransport "questmaster-core/internal/campaign/transport/http"
	characterApp "questmaster-core/internal/character/app"
	characterUsecases "questmaster-core/internal/character/app/usecases"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/middleware"
	userDomain "questmaster-core/internal/user/domain"
)

type fakePremadeRepo struct {
	characterApp.CharacterRepository
}

func (fakePremadeRepo) CreatePremade(name characterDomain.CharacterName, _ campaignDomain.CampaignID, system rpgDomain.System, hp *characterDomain.HP) (characterDomain.Character, error) {
	return characterDomain.Character{Name: name, Slug: "harvey", System: system, Hp: hp}, nil
}

type fakeCampaignFinder struct{ dm userDomain.UserID }

func (f fakeCampaignFinder) FindByID(id campaignDomain.CampaignID) (campaignDomain.Campaign, error) {
	return campaignDomain.Campaign{Id: id, Dm: f.dm, Status: campaignDomain.StatusDraft, System: rpgDomain.CallOfCthulhu}, nil
}

func TestCreatePremadeCharacterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	dm := userDomain.NewUserID(uuid.New())
	handler := NewCharactersHandler(nil, nil, nil, nil, nil, nil,
		characterUsecases.NewCreatePremadeCharacter(fakePremadeRepo{}, fakeCampaignFinder{dm: dm}), nil)

	serve := func(user userDomain.UserID, body string) int {
		router := gin.New()
		router.Use(middleware.ErrorHandlerMiddleware(), func(c *gin.Context) {
			(&context.AppContext{Context: c}).SetUser(userDomain.User{Id: user})
		})
		router.POST("/campaign/:campaignID/character", campaignTransport.CampaignIDMiddleware(), context.Adapt(handler.CreatePremadeCharacter))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/campaign/1/character", strings.NewReader(body)))
		return rec.Code
	}

	if code := serve(dm, `{"name":"Harvey","max_hp":12}`); code != http.StatusCreated {
		t.Fatalf("DM: expected 201, got %d", code)
	}
	if code := serve(userDomain.NewUserID(uuid.New()), `{"name":"Harvey"}`); code != http.StatusForbidden {
		t.Fatalf("non-DM: expected 403, got %d", code)
	}
	for _, body := range []string{`{"name":"  "}`, `{"name":"Harvey","max_hp":0}`, `not json`} {
		if code := serve(dm, body); code != http.StatusBadRequest {
			t.Fatalf("body %s: expected 400, got %d", body, code)
		}
	}
}

type fakeLeaveRepo struct {
	characterApp.CharacterRepository
	player userDomain.UserID
}

func (f fakeLeaveRepo) FindByID(characterDomain.CharacterID) (*characterDomain.Character, error) {
	campaignID := campaignDomain.NewCampaignID(1)
	return &characterDomain.Character{Id: 5, PlayerID: &f.player, CampaignID: &campaignID}, nil
}

func (fakeLeaveRepo) UnlinkFromCampaign(characterDomain.CharacterID, campaignDomain.CampaignID) (bool, error) {
	return true, nil
}

func TestRemoveFromCampaignRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	dm, player := userDomain.NewUserID(uuid.New()), userDomain.NewUserID(uuid.New())
	handler := NewCharactersHandler(nil, nil, nil, nil, nil, nil, nil,
		characterUsecases.NewLeaveCampaign(fakeLeaveRepo{player: player}, fakeCampaignFinder{dm: dm}))

	serve := func(user userDomain.UserID) int {
		router := gin.New()
		router.Use(middleware.ErrorHandlerMiddleware(), func(c *gin.Context) {
			(&context.AppContext{Context: c}).SetUser(userDomain.User{Id: user})
		})
		router.DELETE("/campaign/:campaignID/character/:characterID",
			campaignTransport.CampaignIDMiddleware(), CharacterIDMiddleware(), context.Adapt(handler.RemoveFromCampaign))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/campaign/1/character/5", nil))
		return rec.Code
	}

	if code := serve(player); code != http.StatusNoContent {
		t.Fatalf("player: expected 204, got %d", code)
	}
	if code := serve(userDomain.NewUserID(uuid.New())); code != http.StatusForbidden {
		t.Fatalf("other user: expected 403, got %d", code)
	}
}
