package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"errors": c.Errors.JSON()})
        }
    }
}
