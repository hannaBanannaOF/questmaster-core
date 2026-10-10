package campaign

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/middleware"
)

// serveFilters runs a request through the filters middleware and returns the status and parsed filters
func serveFilters(t *testing.T, mw gin.HandlerFunc, query string) (int, campaignDomain.CampaignListFilters) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	var filters campaignDomain.CampaignListFilters
	router := gin.New()
	router.Use(middleware.ErrorHandlerMiddleware())
	router.GET("/campaign", mw, func(c *gin.Context) {
		filters = (&context.AppContext{Context: c}).CampaignListFilters()
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/campaign"+query, nil))
	return rec.Code, filters
}

func TestCampaignListFiltersMiddleware(t *testing.T) {
	t.Run("no filters", func(t *testing.T) {
		code, filters := serveFilters(t, CampaignListFiltersMiddleware(), "")
		if code != http.StatusOK || filters.Role != nil || filters.Status != nil {
			t.Fatalf("expected 200 without filters, got %d %+v", code, filters)
		}
	})

	t.Run("role and status", func(t *testing.T) {
		code, filters := serveFilters(t, CampaignListFiltersMiddleware(), "?role=player&status=ACTIVE")
		if code != http.StatusOK || filters.Role == nil || *filters.Role != campaignDomain.RolePlayer ||
			filters.Status == nil || *filters.Status != campaignDomain.StatusActive {
			t.Fatalf("expected player and ACTIVE, got %d %+v", code, filters)
		}
	})

	for _, query := range []string{"?role=admin", "?status=FINISHED"} {
		t.Run("unknown value "+query, func(t *testing.T) {
			if code, _ := serveFilters(t, CampaignListFiltersMiddleware(), query); code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d", code)
			}
		})
	}
}

func TestCampaignRoleFilterMiddlewareIgnoresStatus(t *testing.T) {
	code, filters := serveFilters(t, CampaignRoleFilterMiddleware(), "?role=dm&status=FINISHED")
	if code != http.StatusOK || filters.Role == nil || *filters.Role != campaignDomain.RoleDM || filters.Status != nil {
		t.Fatalf("expected role dm and no status, got %d %+v", code, filters)
	}
}
