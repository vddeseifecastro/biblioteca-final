package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		_, _, role := GetCurrentUserFull(c)

		if role != "admin" {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"error": "No tienes permisos para acceder a esta sección",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
