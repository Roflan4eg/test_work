package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			switch {
			case strings.Contains(err.Error(), "not found"):
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
					"code":  "not_found",
				})
			case strings.Contains(err.Error(), "database"):
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
					"code":  "internal_error",
				})
				log.Printf("Database error: %v", err)
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
					"code":  "bad_request",
				})
			}
		}
	}
}
