# 📦 Dependencias del Proyecto

> **Nota:** Este proyecto está desarrollado en **Go**, que gestiona sus dependencias mediante `go.mod` y `go.sum` — el equivalente a `requirements.txt` en Python o `package.json` en Node.js.
> Para instalar todas las dependencias automáticamente, ejecuta `go mod tidy`.

---

## Dependencias principales

| Paquete | Versión | Uso |
|--------|---------|-----|
| `github.com/gin-gonic/gin` | v1.11.0 | Framework web HTTP para el servidor y el router |
| `github.com/golang-jwt/jwt/v5` | v5.3.0 | Generación y validación de tokens JWT para la autenticación |
| `golang.org/x/crypto` | v0.45.0 | Hashing de contraseñas con bcrypt |
| `gorm.io/gorm` | v1.31.1 | ORM para Go — gestión de modelos y consultas a la base de datos |
| `gorm.io/driver/sqlite` | v1.6.0 | Driver SQLite para GORM |

---

## Dependencias indirectas (gestionadas automáticamente por Go)

Estas son instaladas automáticamente como dependencias de las anteriores. No es necesario instalarlas manualmente.

| Paquete | Versión | Razón |
|--------|---------|-------|
| `github.com/bytedance/sonic` | v1.14.2 | Serialización JSON de alto rendimiento (usado por Gin) |
| `github.com/gin-contrib/sse` | v1.1.0 | Soporte de Server-Sent Events para Gin |
| `github.com/go-playground/validator/v10` | v10.28.0 | Validación de structs (usado por Gin) |
| `github.com/mattn/go-sqlite3` | v1.14.32 | Bindings CGO para SQLite (driver nativo) |
| `github.com/jinzhu/inflection` | v1.0.0 | Pluralización de nombres de tablas en GORM |
| `github.com/ugorji/go/codec` | v1.3.1 | Codificación de datos para Gin |
| `google.golang.org/protobuf` | v1.36.10 | Protocol Buffers (dependencia transitiva) |
| `golang.org/x/net` | v0.47.0 | Utilidades de red extendidas |
| `golang.org/x/sys` | v0.38.0 | Acceso a APIs del sistema operativo |
| `golang.org/x/text` | v0.31.0 | Soporte Unicode y normalización de texto |

---

## Instalación

```bash
# Instala todas las dependencias automáticamente
go mod tidy

# Verifica que go.sum esté actualizado y descarga los módulos
go mod download
```

> El archivo `go.sum` contiene los hashes criptográficos de cada dependencia para garantizar integridad en cada instalación.
