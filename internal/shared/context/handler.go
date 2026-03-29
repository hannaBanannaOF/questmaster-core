package context

import "github.com/gin-gonic/gin"

type AppHandler func(*AppContext) error

func Adapt(h AppHandler) gin.HandlerFunc {
	return func(c *gin.Context) {

		appCtx := &AppContext{
			Context: c,
		}

		if err := h(appCtx); err != nil {
			_ = c.Error(err)
			return
		}
	}
}
