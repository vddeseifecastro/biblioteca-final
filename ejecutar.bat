@echo off
echo ========================================
echo   Sistema de Gestión de Biblioteca
echo ========================================
echo.

REM Verificar que Go esté instalado
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: Go no está instalado o no está en el PATH.
    echo Instala Go desde: https://golang.org/dl/
    pause
    exit /b 1
)

REM Verificar estructura de directorios
if not exist "templates\" mkdir templates
if not exist "static\css\" mkdir static\css
if not exist "static\js\" mkdir static\js
if not exist "database\" mkdir database
if not exist "internal\handlers\" mkdir internal\handlers
if not exist "internal\models\" mkdir internal\models
if not exist "internal\database\" mkdir internal\database
if not exist "cmd\web\" mkdir cmd\web

REM Limpiar compilaciones anteriores
echo Limpiando compilaciones anteriores...
if exist "biblioteca-final.exe" del biblioteca-final.exe

REM Instalar dependencias
echo Instalando dependencias...
go mod tidy
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
go get -u golang.org/x/crypto/bcrypt

REM Compilar el proyecto
echo Compilando proyecto...
go build -o biblioteca-final.exe ./cmd/web/

if %errorlevel% neq 0 (
    echo ERROR: Fallo en la compilación.
    pause
    exit /b 1
)

REM Iniciar la aplicación
echo.
echo ========================================
echo   Servidor iniciado correctamente!
echo ========================================
echo.
echo Accede a la aplicación en:
echo http://localhost:8080
echo.
echo Credenciales de prueba:
echo Usuario: admin
echo Contraseña: admin123
echo.
echo Presiona Ctrl+C para detener el servidor.
echo.

REM Ejecutar la aplicación
biblioteca-final.exe