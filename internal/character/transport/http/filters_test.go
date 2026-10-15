package character

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"questmaster-core/internal/shared/middleware"
)

func TestCharacterListSearchFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serve := func(query string) int {
		router := gin.New()
		router.Use(middleware.ErrorHandlerMiddleware())
		router.GET("/character", CharacterListFiltersMiddleware(), func(c *gin.Context) { c.Status(http.StatusOK) })
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/character"+query, nil))
		return rec.Code
	}

	if code := serve("?q=%20%20"); code != http.StatusOK {
		t.Fatalf("blank q: expected 200, got %d", code)
	}
	if code := serve("?q=" + strings.Repeat("a", 101)); code != http.StatusBadRequest {
		t.Fatalf("q too long: expected 400, got %d", code)
	}
}
