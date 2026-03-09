package handlers

import (
	"biblioteca-final/internal/database"
	"biblioteca-final/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ----------------------------------------------------------------------
// VISTAS DE LOGIN Y REGISTRO
// ----------------------------------------------------------------------

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

// ----------------------------------------------------------------------
// REGISTRO
// ----------------------------------------------------------------------

func Register(c *gin.Context) {
	var user models.User

	// Obtener datos del formulario
	user.Name = c.PostForm("name")
	user.Username = c.PostForm("username")
	user.Email = c.PostForm("email")
	password := c.PostForm("password")

	if user.Name == "" || user.Username == "" || user.Email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Todos los campos son obligatorios",
		})
		return
	}

	// Hashear contraseña
	if err := user.HashPassword(password); err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Error al crear usuario",
		})
		return
	}

	// Por defecto, registros públicos siempre son "user"
	role := "user"
	// Si existe un usuario en contexto (registro hecho por admin) y se envió role, permitirlo
	currentUserInterface, exists := c.Get("user")
	if exists && c.PostForm("role") != "" {
		// c.Get("user") en middleware guarda un map[string]interface{} con keys: user_id, username, name, role
		if uMap, ok := currentUserInterface.(map[string]interface{}); ok {
			if r, ok2 := uMap["role"].(string); ok2 && r == "admin" {
				// Admin puede asignar rol desde formulario
				role = c.PostForm("role")
			}
		}
	}
	user.Role = role

	// Guardar usuario
	if err := database.GetDB().Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "El usuario o email ya existen",
		})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

// ----------------------------------------------------------------------
// LOGIN
// ----------------------------------------------------------------------

func Login(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	if err := database.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Usuario o contraseña incorrectos",
		})
		return
	}

	if !user.CheckPassword(password) {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Usuario o contraseña incorrectos",
		})
		return
	}

	// Verificar si el usuario está bloqueado
	if user.BlockedUntil != nil && user.BlockedUntil.After(time.Now()) {
		c.HTML(http.StatusForbidden, "login.html", gin.H{
			"error": "Tu cuenta está bloqueada hasta " + user.BlockedUntil.Format("02/01/2006 15:04"),
		})
		return
	}

	// Crear token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"name":     user.Name,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(getSecretKey()))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": "Error al generar token",
		})
		return
	}

	// Establecer cookie
	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)

	// REDIRECCIÓN SEGÚN ROL
	if user.Role == "admin" {
		c.Redirect(http.StatusFound, "/dashboard/admin")
	} else {
		c.Redirect(http.StatusFound, "/dashboard/user")
	}
}

// ----------------------------------------------------------------------
// LOGOUT
// ----------------------------------------------------------------------

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
