package campaign

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/middleware"
	userDomain "questmaster-core/internal/user/domain"
)

type fakeCountRepo struct {
	campaignApp.CampaignRepository
}

func (fakeCountRepo) CountByStatusForUser(userDomain.UserID, campaignDomain.CampaignListFilters) (map[campaignDomain.CampaignStatus]int, error) {
	return map[campaignDomain.CampaignStatus]int{campaignDomain.StatusActive: 2}, nil
}

// countsRouter registers the counts route next to /:campaignID, like cmd/app/routes does
func countsRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	handler := NewCampaignHandler(nil, nil, nil, nil, nil, nil, campaignUsecases.NewGetCampaignStatusCounts(fakeCountRepo{}), nil)

	router := gin.New()
	router.Use(middleware.ErrorHandlerMiddleware(), func(c *gin.Context) {
		(&context.AppContext{Context: c}).SetUser(userDomain.User{Id: userDomain.NewUserID(uuid.New())})
	})
	router.GET("/campaign/counts", CampaignRoleFilterMiddleware(), context.Adapt(handler.GetStatusCounts))
	router.GET("/campaign/:campaignID", CampaignIDMiddleware(), func(c *gin.Context) { c.Status(http.StatusTeapot) })
	return router
}

func TestGetStatusCountsRoute(t *testing.T) {
	t.Run("returns every status", func(t *testing.T) {
		rec := httptest.NewRecorder()
		countsRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/campaign/counts?role=dm", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body)
		}
		var body map[string]int
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(body) != 4 || body["ACTIVE"] != 2 || body["DRAFT"] != 0 || body["PAUSED"] != 0 || body["ARCHIVED"] != 0 {
			t.Fatalf("unexpected body %v", body)
		}
	})

	t.Run("unknown role", func(t *testing.T) {
		rec := httptest.NewRecorder()
		countsRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/campaign/counts?role=admin", nil))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}

type fakeCreateRepo struct {
	campaignApp.CampaignRepository
	dmName *rpgDomain.TableName
}

func (f *fakeCreateRepo) Create(name campaignDomain.CampaignName, _ *campaignDomain.CampaignOverview, dm userDomain.UserID, system rpgDomain.System, dmName rpgDomain.TableName) (campaignDomain.Campaign, error) {
	f.dmName = &dmName
	return campaignDomain.Campaign{Id: 1, Name: name, Dm: dm, System: system, DmName: dmName}, nil
}

func TestCreateCampaignRoute(t *testing.T) {
	create := func(body string) (*httptest.ResponseRecorder, *fakeCreateRepo) {
		gin.SetMode(gin.TestMode)
		repo := &fakeCreateRepo{}
		handler := NewCampaignHandler(nil, nil, campaignUsecases.NewCreateCampaign(repo), nil, nil, nil, nil, nil)
		router := gin.New()
		router.Use(middleware.ErrorHandlerMiddleware(), func(c *gin.Context) {
			(&context.AppContext{Context: c}).SetUser(userDomain.User{Id: userDomain.NewUserID(uuid.New())})
		})
		router.POST("/campaign", context.Adapt(handler.CreateCampaign))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/campaign", strings.NewReader(body)))
		return rec, repo
	}

	t.Run("stores the trimmed DM name", func(t *testing.T) {
		rec, repo := create(`{"name":"Curse","system":"DUNGEONS_AND_DRAGONS","dm_name":" Mestre Ana "}`)
		if rec.Code != http.StatusCreated || repo.dmName == nil || *repo.dmName != "Mestre Ana" {
			t.Fatalf("expected 201 with the DM name stored, got %d %s saved=%v", rec.Code, rec.Body, repo.dmName)
		}
	})

	for label, body := range map[string]string{
		"missing": `{"name":"Curse","system":"DUNGEONS_AND_DRAGONS"}`,
		"blank":   `{"name":"Curse","system":"DUNGEONS_AND_DRAGONS","dm_name":"   "}`,
	} {
		t.Run(label+" DM name", func(t *testing.T) {
			rec, repo := create(body)
			if rec.Code != http.StatusBadRequest || repo.dmName != nil {
				t.Fatalf("expected 400 without creating, got %d saved=%v", rec.Code, repo.dmName)
			}
		})
	}
}
