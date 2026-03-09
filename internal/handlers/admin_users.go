package handlers

import (
	"biblioteca-final/internal/database"
	"biblioteca-final/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminUsers muestra la vista de administración de usuarios
func AdminUsers(c *gin.Context) {
	db := database.GetDB()

	var users []models.User
	db.Where("deleted_at IS NULL").Find(&users) // solo usuarios no eliminados

	// Contadores
	totalUsers := len(users)
	blockedUsers := 0
	adminUsers := 0

	for _, u := range users {
		if u.BlockedUntil != nil {
			blockedUsers++
		}
		if u.Role == "admin" {
			adminUsers++
		}
	}

	// Nombre del admin actual
	name, _, _ := GetCurrentUserFull(c)

	c.HTML(http.StatusOK, "users_admin.html", gin.H{
		"users":        users,
		"totalUsers":   totalUsers,
		"blockedUsers": blockedUsers,
		"adminUsers":   adminUsers,
		"name":         name,
	})
}

// BlockUser bloquea un usuario por 14 días
func BlockUser(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Usuario no encontrado"})
		return
	}

	duration := time.Hour * 24 * 14
	until := time.Now().Add(duration)
	user.BlockedUntil = &until
	db.Save(&user)

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

// UnblockUser desbloquea un usuario
func UnblockUser(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Usuario no encontrado"})
		return
	}

	user.BlockedUntil = nil
	db.Save(&user)

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

// DeleteUser hace un soft delete del usuario
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	currentID, _, _ := GetCurrentUserFull(c)
	if id == fmt.Sprint(currentID) {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"error": "No puedes eliminar tu propio usuario",
		})
		return
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Soft delete: en lugar de db.Delete(&user), usamos gorm.Delete
	if err := db.Delete(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "No se pudo eliminar el usuario"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/users")
}

// ClearLoanReports elimina todos los reportes históricos
func ClearLoanReports(c *gin.Context) {
	db := database.GetDB()

	if err := db.Exec("DELETE FROM loan_reports").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron eliminar los reportes"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/loans/admin")
}

// ClearReturnedLoans elimina solo los préstamos devueltos
func ClearReturnedLoans(c *gin.Context) {
	db := database.GetDB()

	if err := db.Where("returned = ?", true).Delete(&models.Loan{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron eliminar los préstamos devueltos"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/loans/admin")
}
