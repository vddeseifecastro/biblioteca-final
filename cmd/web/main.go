package main

import (
	"biblioteca-final/internal/database"
	"biblioteca-final/internal/handlers"
	"biblioteca-final/internal/models"
	"bufio"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Crear usuario administrador por defecto si no existe
func createDefaultUser() {
	var user models.User
	if err := database.GetDB().Where("username = ?", "admin").First(&user).Error; err != nil {
		adminUser := models.User{
			Username: "admin",
			Email:    "admin@biblioteca.com",
			Name:     "Administrador",
			Role:     "admin",
		}
		adminUser.HashPassword("admin123")
		database.GetDB().Create(&adminUser)
		log.Println("✅ Usuario administrador creado: admin / admin123")
	} else {
		log.Println("ℹ️ Usuario administrador ya existe")
	}
}

// loadEnv lee el archivo .env y carga las variables de entorno
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return // si no existe .env, usa las variables del sistema
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
}

func main() {
	// Cargar variables de entorno desde .env
	loadEnv()

	// Conectar a la base de datos
	if err := database.Connect(); err != nil {
		log.Fatal("❌ Error conectando a la base de datos:", err)
	}

	// Migrar tablas
	database.Migrate()

	// Crear usuario admin por defecto
	createDefaultUser()

	// Modo producción (elimina warnings del navegador)
	gin.SetMode(gin.ReleaseMode)

	// Configurar router Gin
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Funciones para templates
	router.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a % b
		},
		"gt":  func(a, b int) bool { return a > b },
		"lt":  func(a, b int) bool { return a < b },
		"eq":  func(a, b interface{}) bool { return a == b },
		"ne":  func(a, b interface{}) bool { return a != b },
		"now": func() time.Time { return time.Now() },
	})

	// Cargar templates y estáticos
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static")

	// ---------------------
	// RUTAS PÚBLICAS
	// ---------------------
	router.GET("/", handlers.ShowLoginPage)
	router.GET("/login", handlers.ShowLoginPage)
	router.POST("/login", handlers.Login)
	router.GET("/register", handlers.ShowRegisterPage)
	router.POST("/register", handlers.Register)

	// ---------------------
	// RUTAS PROTEGIDAS
	// ---------------------
	protected := router.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.GET("/dashboard", handlers.RedirectDashboardByRole)
		protected.GET("/dashboard/admin", handlers.ShowDashboardAdmin)
		protected.GET("/dashboard/user", handlers.ShowDashboardUser)
		protected.GET("/books", handlers.ListBooks)
		protected.GET("/books/:id", handlers.ShowBook)
		protected.GET("/books/:id/borrow", handlers.BorrowBook)
		protected.GET("/loans", handlers.ListLoansUser)
		protected.GET("/logout", handlers.Logout)
		protected.GET("/profile", handlers.ShowUserProfile)

		// ---------------------
		// RUTAS ADMIN
		// ---------------------
		admin := protected.Group("/")
		admin.Use(handlers.AdminMiddleware())
		{
			admin.GET("/books/new", handlers.ShowCreateBook)
			admin.POST("/books", handlers.CreateBook)
			admin.GET("/books/:id/edit", handlers.ShowEditBook)
			admin.POST("/books/:id", handlers.UpdateBook)
			admin.POST("/books/:id/delete", handlers.DeleteBook)
			// Ruta para actualizar stock
			admin.POST("/books/:id/stock", handlers.UpdateBookStock)

			// Préstamos admin
			admin.GET("/loans/admin", handlers.ListLoansAdmin)
			admin.POST("/loans/:id/return", handlers.ReturnBook)
			//Ruta eliminar informes
			admin.POST("/loans/admin/clear-reports", handlers.ClearLoanReports)
			admin.POST("/loans/admin/clear-returned", handlers.ClearReturnedLoans)
			admin.GET("/admin/loans/report", handlers.GenerateLoanReport)

			// Gestión de usuarios admin
			admin.GET("/admin/users", handlers.AdminUsers)
			admin.POST("/admin/users/:id/block", handlers.BlockUser)
			admin.POST("/admin/users/:id/unblock", handlers.UnblockUser)
			admin.POST("/admin/users/:id/delete", handlers.DeleteUser)
		}
	}

	// Información de inicio
	log.Println("========================================")
	log.Println("   SISTEMA DE GESTIÓN DE BIBLIOTECA")
	log.Println("========================================")
	log.Println("🌐 Servidor: http://localhost:8080")
	log.Println("========================================")

	// Iniciar servidor
	if err := router.Run(":8080"); err != nil {
		log.Fatal("❌ Error iniciando servidor:", err)
	}
}
