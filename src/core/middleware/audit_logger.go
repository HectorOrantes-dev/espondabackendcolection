package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func AuditLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		userID, _ := c.Get("userID")
		email, _ := c.Get("email")

		log.Printf("[AUDIT] %s | %d | %s | %s | user=%v | email=%v",
			time.Now().Format(time.RFC3339),
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			userID,
			email,
		)
		_ = start
	}
}
