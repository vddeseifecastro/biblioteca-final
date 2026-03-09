Write-Host "========================================" -ForegroundColor Green
Write-Host "  Sistema de Gestión de Biblioteca" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Green
Write-Host ""

# Verificar que Go esté instalado
$goVersion = go version 2>$null
if (-not $goVersion) {
    Write-Host "ERROR: Go no está instalado o no está en el PATH." -ForegroundColor Red
    Write-Host "Instala Go desde: https://golang.org/dl/" -ForegroundColor Yellow
    pause
    exit 1
}

Write-Host "Go versión: $goVersion" -ForegroundColor Cyan

# Crear directorios necesarios
$directories = @(
    "templates",
    "static\css", 
    "static\js",
    "static\images",
    "database",
    "internal\handlers",
    "internal\models",
    "internal\database",
    "cmd\web"
)

foreach ($dir in $directories) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
        Write-Host "Creado directorio: $dir" -ForegroundColor DarkGray
    }
}

# Limpiar compilaciones anteriores
if (Test-Path "biblioteca-final.exe") {
    Remove-Item "biblioteca-final.exe" -Force
    Write-Host "Compilaciones anteriores eliminadas" -ForegroundColor DarkGray
}

# Inicializar módulo si no existe
if (-not (Test-Path "go.mod")) {
    Write-Host "Inicializando módulo Go..." -ForegroundColor Cyan
    go mod init biblioteca-final
}

# Instalar/actualizar dependencias
Write-Host "Instalando dependencias..." -ForegroundColor Cyan
go mod tidy
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
go get -u golang.org/x/crypto/bcrypt

# Ejecutar la aplicación
Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  Iniciando servidor..." -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Accede a la aplicación en:" -ForegroundColor White
Write-Host "  http://localhost:8080" -ForegroundColor Cyan -BackgroundColor DarkGray
Write-Host ""
Write-Host "Credenciales de prueba:" -ForegroundColor White
Write-Host "  Usuario: admin" -ForegroundColor Yellow
Write-Host "  Contraseña: admin123" -ForegroundColor Yellow
Write-Host ""
Write-Host "Para crear un nuevo usuario, usa /register" -ForegroundColor DarkGray
Write-Host "Presiona Ctrl+C para detener el servidor." -ForegroundColor DarkGray
Write-Host ""

go run cmd/web/main.go