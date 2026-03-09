package handlers

import (
	"biblioteca-final/internal/database"
	"biblioteca-final/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Perfil de usuario: muestra datos, historial y estadísticas
func ShowUserProfile(c *gin.Context) {
	username, name := GetCurrentUser(c)
	if username == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := database.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Usuario no encontrado"})
		return
	}

	var loans []models.Loan
	database.GetDB().Preload("Book").Where("user_id = ?", user.ID).Find(&loans)

	// Estadísticas
	totalLoans := len(loans)
	returned := 0
	late := 0
	for _, loan := range loans {
		if loan.Returned {
			returned++
		} else if loan.DueDate != nil && loan.DueDate.Before(time.Now()) {
			late++
		}
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"username":    username,
		"name":        name,
		"email":       user.Email,
		"role":        user.Role,
		"blocked":     user.BlockedUntil,
		"preferences": user.Preferences,
		"loans":       loans,
		"totalLoans":  totalLoans,
		"returned":    returned,
		"late":        late,
	})
}
