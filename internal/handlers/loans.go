package handlers

import (
	"biblioteca-final/internal/database"
	"biblioteca-final/internal/models"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ListLoansUser muestra los préstamos del usuario logueado
func ListLoansUser(c *gin.Context) {
	username, _ := GetCurrentUser(c)
	user := GetUserByUsername(username)

	var loans []models.Loan
	database.GetDB().Preload("Book").Where("user_id = ?", user.ID).Find(&loans)

	c.HTML(http.StatusOK, "loans_user.html", gin.H{
		"username": username,
		"name":     user.Name,
		"loans":    loans,
	})
}

// ListLoansAdmin muestra todos los préstamos (solo admin)
func ListLoansAdmin(c *gin.Context) {
	db := database.GetDB()
	var loans []models.Loan
	db.Preload("User").Preload("Book").Find(&loans)

	// Si se envía parámetro clear=true, vaciamos la vista
	if c.Query("clear") == "true" {
		loans = []models.Loan{}
	}

	c.HTML(http.StatusOK, "loans_admin.html", gin.H{
		"name":     "Admin",
		"username": "admin",
		"loans":    loans,
	})
}

func GenerateLoanReport(c *gin.Context) {
	db := database.GetDB()

	var loans []models.Loan
	db.Preload("User").Preload("Book").Find(&loans)

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=reporte_prestamos.csv")
	c.Header("Cache-Control", "no-cache")

	writer := csv.NewWriter(c.Writer)

	// Header del CSV
	writer.Write([]string{
		"ID",
		"Usuario",
		"Libro",
		"Fecha Prestamo",
		"Fecha Limite",
		"Devuelto",
		"Estado",
	})

	for _, loan := range loans {
		returned := "No"
		if loan.Returned {
			returned = "Sí"
		}

		state := "Activo"
		if loan.Returned {
			state = "Completado"
		} else if loan.DueDate != nil && loan.DueDate.Before(time.Now()) {
			state = "Vencido"
		}

		writer.Write([]string{
			fmt.Sprintf("%d", loan.ID),
			loan.User.Name,
			loan.Book.Title,
			loan.CreatedAt.Format("02/01/2006"),
			func() string {
				if loan.DueDate != nil {
					return loan.DueDate.Format("02/01/2006")
				}
				return "-"
			}(),
			returned,
			state,
		})
	}

	writer.Flush()
}
