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
		characterUsecases.NewCreatePremadeCharacter(fakePremadeRepo{}, fakeCampaignFinder{dm: dm}))

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
