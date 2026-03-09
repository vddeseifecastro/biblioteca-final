package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() string {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		panic("JWT_SECRET_KEY no está definida en las variables de entorno")
	}
	return key
}

// ----------------------------------------------------------------------
// 🔐 MIDDLEWARE DE AUTENTICACIÓN GENERAL
// ----------------------------------------------------------------------

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(getSecretKey()), nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Guardar usuario en contexto
		c.Set("user", map[string]interface{}{
			"user_id":  claims["user_id"],
			"username": claims["username"],
			"name":     claims["name"],
			"role":     claims["role"],
		})

		c.Next()
	}
}

// ----------------------------------------------------------------------
// 🧩 FUNCIONES PARA OBTENER USUARIO ACTUAL
// ----------------------------------------------------------------------

func GetCurrentUser(c *gin.Context) (string, string) {
	user, exists := c.Get("user")
	if !exists {
		return "", ""
	}

	u := user.(map[string]interface{})
	return u["username"].(string), u["name"].(string)
}

func GetCurrentUserFull(c *gin.Context) (string, string, string) {
	user, exists := c.Get("user")
	if !exists {
		return "", "", ""
	}

	u := user.(map[string]interface{})
	return u["username"].(string), u["name"].(string), u["role"].(string)
}
