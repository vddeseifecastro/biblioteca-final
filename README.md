<div align="center">

# 📚 BiblioSystem
### Full Stack Library Management System

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.11-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-1.31-FF6B6B?style=for-the-badge)
![Bootstrap](https://img.shields.io/badge/Bootstrap-5-7952B3?style=for-the-badge&logo=bootstrap&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)

A full-featured library management system built from scratch with Go, Gin and GORM. Includes JWT authentication, book catalog, loan management, admin panel and user administration.

[Report Bug](https://github.com/vddeseifecastro/biblioteca-final/issues) · [Request Feature](https://github.com/vddeseifecastro/biblioteca-final/issues)

</div>

---

## 📸 Screenshots

### 🔐 Authentication

**Login**

![Login](https://github.com/user-attachments/assets/8f5c1aae-468f-4475-aabe-37d83af11e1c)

**Register**

![Register](https://github.com/user-attachments/assets/685a56e7-d25d-49ed-a897-ea86ccac0238)

### ⚙️ Admin Dashboard

![Admin Dashboard 1](https://github.com/user-attachments/assets/ad4d8a69-832c-404c-8032-baf72e391563)

![Admin Dashboard 2](https://github.com/user-attachments/assets/988da001-3c95-4264-a262-24ca764aae56)

### 📚 Books

![Books 1](https://github.com/user-attachments/assets/b6d66b0d-9887-4493-85bd-70c03cbc025c)

![Books 2](https://github.com/user-attachments/assets/8ddcea02-1e55-4d88-9840-cfea6f671b21)

### ➕ New Book

![New Book](https://github.com/user-attachments/assets/a0aea398-52f9-4f33-96f2-f09ba67ceecf)

### 👥 User Administration

![User Administration](https://github.com/user-attachments/assets/366c26b6-349e-42f3-ad8a-788df64f48a9)

### 🔍 Book Detail

![Book Detail](https://github.com/user-attachments/assets/11016ff3-fcb1-484f-b33f-e3a23c5866e8)

### ✏️ Edit Book

![Edit Book](https://github.com/user-attachments/assets/d3c6f48a-016c-446f-98a6-a774aff38047)

### 🔄 Loan Management

![Loan Management](https://github.com/user-attachments/assets/c67d9eee-574c-41bd-80e2-9d6d5112f253)

### 🏠 User Panel

![User Panel](https://github.com/user-attachments/assets/f0098c4e-5c00-45da-8a77-8a4bb712127d)

### 🚫 Unauthorized Access Error

![Unauthorized Error](https://github.com/user-attachments/assets/04ee84c3-8223-447f-9dd9-794e27d467b4)

---

## ✨ Features

### 👤 User
- JWT authentication (register & login)
- Book catalog with search by title, author or ISBN
- Filter by category
- Book detail page
- Book loans with automatic due date (14 days)
- Loan history with real-time status
- Personal profile with reading statistics

### 🛠️ Admin Panel
- Full book CRUD with cover image upload
- Real-time stock management
- Global view of all active, overdue and returned loans
- User management (block, unblock, delete)
- CSV report export
- Dashboard with system-wide statistics

---

## 🖥️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25, Gin v1.11 |
| ORM | GORM v1.31 + SQLite |
| Auth | JWT (`golang-jwt/jwt v5`) + bcrypt |
| Frontend | HTML, Bootstrap 5, Custom CSS |
| Security | Environment variables, HttpOnly cookies |

---

## 🚀 Getting Started

### Prerequisites
- Go 1.19+
- Git

### Setup

```bash
# 1. Clone the repository
git clone https://github.com/vddeseifecastro/biblioteca-final.git
cd biblioteca-final

# 2. Configure environment variables
cp .env.example .env
# Edit .env and add a secure JWT key:
# JWT_SECRET_KEY=your_secret_key_here
# You can generate one with: openssl rand -hex 32

# 3. Download dependencies
go mod tidy

# 4. Run the server
go run ./cmd/web/main.go
```

Server running at `http://localhost:8080`

Default admin credentials:
- **Username:** `admin`
- **Password:** `admin123`

---

## 📁 Project Structure

```
biblioteca-final/
├── cmd/
│   └── web/
│       └── main.go              # Entry point, routes and server config
├── internal/
│   ├── database/
│   │   └── database.go          # SQLite connection and GORM migrations
│   ├── handlers/
│   │   ├── auth.go              # Login, register and logout
│   │   ├── middleware.go        # AuthMiddleware and AdminMiddleware (JWT)
│   │   ├── books.go             # Book CRUD and loans
│   │   ├── loans.go             # Loan management and CSV export
│   │   ├── admin.go             # Admin-only middleware
│   │   ├── admin_users.go       # User management (block/delete)
│   │   └── profile.go           # User profile and statistics
│   └── models/
│       ├── user.go              # User model with bcrypt
│       ├── book.go              # Book model
│       ├── loan.go              # Loan model
│       └── loan_report.go       # LoanReport model
├── templates/                   # HTML templates with Go template engine
├── static/                      # CSS, JS and static images
├── .env.example                 # Environment variables template
├── .gitignore
└── go.mod
```

---

## 🔄 Loan Status Flow

```
borrowed → active → overdue
                       ↓
                    returned
```

---

## 🔐 Security

- Passwords stored with **bcrypt** (cost=14)
- JWT tokens signed with a key defined in environment variables (`JWT_SECRET_KEY`)
- Session cookies are **HttpOnly** to prevent JavaScript access
- `.env` file and database excluded from repository via `.gitignore`
- Admin role validated server-side on every protected route

---

## 🌱 Upcoming Features

- [ ] Deploy on Railway or Render
- [ ] Email notifications on loan status change
- [ ] Book reviews & ratings
- [ ] Advanced search filters
- [ ] Pagination for large catalogs

---

## 👨‍💻 Author

**Victor Dominic Deseife Castro**

[![GitHub](https://img.shields.io/badge/GitHub-vddeseifecastro-181717?style=for-the-badge&logo=github)](https://github.com/vddeseifecastro)

---

<div align="center">
  <p>Built with ❤️ by Victor Dominic Deseife Castro</p>
  <p>⭐ Star this repo if you found it useful!</p>
</div>