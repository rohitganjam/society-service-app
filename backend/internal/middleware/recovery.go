package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rohit/society-service-app/backend/internal/utils"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
				utils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
