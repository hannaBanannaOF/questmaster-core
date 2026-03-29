package context

import "github.com/gin-gonic/gin"

type AppContext struct {
	*gin.Context
}
