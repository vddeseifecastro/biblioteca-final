package handlers

import (
	"biblioteca-final/internal/database"
	"biblioteca-final/internal/models"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// -----------------------
// DASHBOARD
// -----------------------
func RedirectDashboardByRole(c *gin.Context) {
	_, _, role := GetCurrentUserFull(c)

	if role == "admin" {
		c.Redirect(http.StatusFound, "/dashboard/admin")
	} else {
		c.Redirect(http.StatusFound, "/dashboard/user")
	}
}

func ShowDashboardAdmin(c *gin.Context) {
	db := database.GetDB()

	var totalBooks int64
	db.Model(&models.Book{}).Count(&totalBooks)

	var availableBooks int64
	db.Model(&models.Book{}).Where("available > ?", 0).Count(&availableBooks)

	var unavailableBooks int64
	db.Model(&models.Book{}).Where("available = ?", 0).Count(&unavailableBooks)

	var totalUsers int64
	db.Model(&models.User{}).Count(&totalUsers)

	var activeLoans int64
	db.Model(&models.Loan{}).Where("returned = ?", false).Count(&activeLoans)

	var overdueLoans int64
	db.Model(&models.Loan{}).Where("returned = ? AND due_date < ?", false, time.Now()).Count(&overdueLoans)

	var recentUsers []models.User
	db.Order("created_at desc").Limit(5).Find(&recentUsers)

	var recentLoans []models.Loan
	db.Preload("User").Preload("Book").Order("created_at desc").Limit(6).Find(&recentLoans)

	name, _, _ := GetCurrentUserFull(c)

	c.HTML(http.StatusOK, "dashboard_admin.html", gin.H{
		"totalBooks":       totalBooks,
		"availableBooks":   availableBooks,
		"unavailableBooks": unavailableBooks,
		"totalUsers":       totalUsers,
		"activeLoans":      activeLoans,
		"overdueLoans":     overdueLoans,
		"recentUsers":      recentUsers,
		"recentLoans":      recentLoans,
		"name":             name,
	})
}

func ShowDashboardUser(c *gin.Context) {
	db := database.GetDB()
	username, name, _ := GetCurrentUserFull(c)

	var user models.User
	db.Where("username = ?", username).First(&user)

	var loans []models.Loan
	db.Preload("Book").Where("user_id = ? AND returned = ?", user.ID, false).Find(&loans)

	now := time.Now()
	for i := range loans {
		if loans[i].DueDate != nil && loans[i].DueDate.Before(now) && !loans[i].Returned {
			loans[i].IsOverdue = true
		}
	}

	var overdue int64
	db.Model(&models.Loan{}).Where("user_id = ? AND returned = ? AND due_date < ?", user.ID, false, time.Now()).Count(&overdue)

	c.HTML(http.StatusOK, "dashboard_user.html", gin.H{
		"name":      name,
		"loans":     loans,
		"overdue":   overdue,
		"loanCount": len(loans),
	})
}

// -----------------------
// LISTADO DE LIBROS — (sin paginación forzada)
// -----------------------
func ListBooks(c *gin.Context) {
	// Obtener usuario + rol
	username, name, role := GetCurrentUserFull(c)

	// Filtros
	search := c.Query("search")
	category := c.Query("category") // si en el futuro quieres filtrar por categoría

	var books []models.Book
	query := database.GetDB().Where("is_active = ?", true)

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title LIKE ? OR author LIKE ? OR isbn LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Traer todos los libros (sin Limit), ordenados por más recientes
	query.Order("created_at ASC").Find(&books)

	// Categorías (para un filtro UI)
	var categories []string
	database.GetDB().Model(&models.Book{}).Distinct().Pluck("category", &categories)

	c.HTML(http.StatusOK, "books.html", gin.H{
		"username":   username,
		"name":       name,
		"role":       role,
		"books":      books,
		"categories": categories,
		"search":     search,
		"category":   category,
	})
}

// -----------------------
// CREAR LIBRO
// -----------------------
func ShowCreateBook(c *gin.Context) {
	username, name, _ := GetCurrentUserFull(c)
	c.HTML(http.StatusOK, "new_book.html", gin.H{
		"username": username,
		"name":     name,
	})
}

func CreateBook(c *gin.Context) {
	username, name, _ := GetCurrentUserFull(c)
	var book models.Book

	publishedAtStr := c.PostForm("publishedAt")
	if publishedAtStr != "" {
		if t, err := time.Parse("2006-01-02", publishedAtStr); err == nil {
			book.PublishedAt = &t
		}
	}

	book.ISBN = c.PostForm("isbn")
	book.Title = c.PostForm("title")
	book.Author = c.PostForm("author")
	book.Description = c.PostForm("description")
	book.Category = c.PostForm("category")
	book.Publisher = c.PostForm("publisher")
	book.Pages, _ = strconv.Atoi(c.PostForm("pages"))
	book.Language = c.PostForm("language")
	book.Stock, _ = strconv.Atoi(c.PostForm("stock"))
	book.Available = book.Stock

	// Subida de imagen (campo coverImage -> archivo)
	file, err := c.FormFile("coverImage")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.HTML(http.StatusBadRequest, "new_book.html", gin.H{
				"error":    "Formato de imagen no válido (solo jpg, jpeg, png)",
				"username": username,
				"name":     name,
			})
			return
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := "./static/covers/" + filename
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.HTML(http.StatusInternalServerError, "new_book.html", gin.H{
				"error":    "Error al guardar la imagen",
				"username": username,
				"name":     name,
			})
			return
		}
		book.CoverImage = "/static/covers/" + filename
	} else {
		// Si no sube archivo dejamos placeholder
		book.CoverImage = "/static/covers/placeholder.png"
	}

	if err := database.GetDB().Create(&book).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "new_book.html", gin.H{
			"error":    "Error al crear el libro (ISBN ya existe)",
			"username": username,
			"name":     name,
		})
		return
	}

	c.Redirect(http.StatusFound, "/books")
}

// -----------------------
// VER UN LIBRO
// -----------------------
func ShowBook(c *gin.Context) {
	username, name, _ := GetCurrentUserFull(c)
	id := c.Param("id")
	var book models.Book
	if err := database.GetDB().First(&book, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error":    "Libro no encontrado",
			"username": username,
			"name":     name,
		})
		return
	}
	c.HTML(http.StatusOK, "view_book.html", gin.H{
		"username": username,
		"name":     name,
		"book":     book,
	})
}

// -----------------------
// EDITAR LIBRO
// -----------------------
func ShowEditBook(c *gin.Context) {
	username, name, _ := GetCurrentUserFull(c)
	id := c.Param("id")
	var book models.Book
	if err := database.GetDB().First(&book, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error":    "Libro no encontrado",
			"username": username,
			"name":     name,
		})
		return
	}
	c.HTML(http.StatusOK, "edit_book.html", gin.H{
		"username": username,
		"name":     name,
		"book":     book,
	})
}

func UpdateBook(c *gin.Context) {
	username, name, _ := GetCurrentUserFull(c)
	id := c.Param("id")
	var book models.Book
	if err := database.GetDB().First(&book, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error":    "Libro no encontrado",
			"username": username,
			"name":     name,
		})
		return
	}

	book.Title = c.PostForm("Title")
	book.Author = c.PostForm("Author")
	book.Description = c.PostForm("Description")
	book.Category = c.PostForm("Category")
	book.Publisher = c.PostForm("Publisher")
	publishedAtStr := c.PostForm("PublishedAt")
	if publishedAtStr != "" {
		if t, err := time.Parse("2006-01-02", publishedAtStr); err == nil {
			book.PublishedAt = &t
		}
	}
	pages, _ := strconv.Atoi(c.PostForm("Pages"))
	book.Pages = pages
	book.Language = c.PostForm("Language")
	stock, err := strconv.Atoi(c.PostForm("Stock"))
	if err != nil || stock < 0 {
		c.HTML(http.StatusBadRequest, "edit_book.html", gin.H{
			"error":    "Stock inválido",
			"book":     book,
			"username": username,
			"name":     name,
		})
		return
	}

	var loansCount int64
	database.GetDB().Model(&models.Loan{}).Where("book_id = ? AND returned = ?", book.ID, false).Count(&loansCount)
	book.Stock = stock
	disp := stock - int(loansCount)
	if disp < 0 {
		disp = 0
	}
	book.Available = disp

	// Subida de nueva portada
	file, err := c.FormFile("coverImage")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.HTML(http.StatusBadRequest, "edit_book.html", gin.H{
				"error":    "Formato de imagen no válido (solo jpg, jpeg, png)",
				"book":     book,
				"username": username,
				"name":     name,
			})
			return
		}
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := "./static/covers/" + filename
		if err := c.SaveUploadedFile(file, savePath); err == nil {
			book.CoverImage = "/static/covers/" + filename
		}
	}

	database.GetDB().Save(&book)
	c.Redirect(http.StatusFound, "/books")
}

// -----------------------
// ELIMINAR LIBRO
// -----------------------
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := database.GetDB().Delete(&models.Book{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el libro"})
		return
	}
	c.Redirect(http.StatusFound, "/books")
}

// -----------------------
// TOMAR PRESTADO
// -----------------------
func BorrowBook(c *gin.Context) {
	username, _ := GetCurrentUser(c)
	user := GetUserByUsername(username)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Usuario no válido"})
		return
	}

	bookID := c.Param("id")
	var book models.Book
	if err := database.GetDB().First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Libro no encontrado"})
		return
	}

	if book.Available <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay copias disponibles"})
		return
	}

	due := time.Now().AddDate(0, 0, 14)
	loan := models.Loan{
		UserID:  user.ID,
		BookID:  book.ID,
		DueDate: &due,
	}

	if err := database.GetDB().Create(&loan).Error; err != nil {
		log.Println("ERROR creando préstamo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el préstamo"})
		return
	}

	book.Available--
	database.GetDB().Save(&book)

	c.Redirect(http.StatusFound, "/books")
}

// -----------------------
// DEVOLVER LIBRO
// -----------------------
func ReturnBook(c *gin.Context) {
	loanID := c.Param("id")
	var loan models.Loan
	if err := database.GetDB().First(&loan, loanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Préstamo no encontrado"})
		return
	}

	if !loan.Returned {
		loan.Returned = true
		var book models.Book
		database.GetDB().First(&book, loan.BookID)
		book.Available++
		database.GetDB().Save(&book)
		database.GetDB().Save(&loan)
	}

	c.Redirect(http.StatusFound, "/loans/admin")
}

// -----------------------
// FUNCIONES AUXILIARES
// -----------------------
func GetUserByUsername(username string) models.User {
	var user models.User
	database.GetDB().Where("username = ?", username).First(&user)
	return user
}

// -----------------------
// ACTUALIZAR STOCK
// -----------------------
func UpdateBookStock(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := database.GetDB().First(&book, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Libro no encontrado"})
		return
	}

	stockStr := c.PostForm("stock")
	stock, err := strconv.Atoi(stockStr)
	if err != nil || stock < 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Stock inválido"})
		return
	}

	book.Stock = stock
	if book.Available > stock {
		book.Available = stock
	}
	if book.Available < 0 {
		book.Available = 0
	}
	database.GetDB().Save(&book)
	c.Redirect(http.StatusFound, "/books")
}
