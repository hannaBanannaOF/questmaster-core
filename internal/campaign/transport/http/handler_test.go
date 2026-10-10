package campaign

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	campaignApp "questmaster-core/internal/campaign/app"
	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignDomain "questmaster-core/internal/campaign/domain"
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
	handler := NewCampaignHandler(nil, nil, nil, nil, nil, nil, campaignUsecases.NewGetCampaignStatusCounts(fakeCountRepo{}))

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
