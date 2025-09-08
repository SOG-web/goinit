# GoInit - Go Gin API Generator

[![Go Report Card](https://goreportcard.com/badge/github.com/rou/goinit)](https://goreportcard.com/report/github.com/rou/goinit)
[![GoDoc](https://godoc.org/github.com/rou/goinit?status.svg)](https://godoc.org/github.com/rou/goinit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful CLI tool to generate production-ready Go API projects built with the Gin framework.

## 🚀 Installation

### Option 1: Go Install (Recommended)

```bash
go install github.com/rou/goinit@latest
```

### Option 2: Download Binary

Download the latest release from [GitHub Releases](https://github.com/rou/goinit/releases)

### Option 3: Build from Source

```bash
git clone https://github.com/rou/goinit.git
cd goinit
go build -o goinit .
```

## 📖 Usage

```bash
goinit
```

Follow the interactive prompts to configure your new Go API project.

## ✨ Features

- **Interactive Setup**: Step-by-step project configuration
- **Multiple Databases**: SQLite, MySQL, PostgreSQL support
- **Authentication**: JWT and session-based auth
- **Real-time**: SSE and WebSocket endpoints
- **Email Integration**: SMTP and local email logging
- **File Storage**: Local and S3-compatible storage
- **Admin Panel**: Built-in admin functionality
- **API Documentation**: Auto-generated Swagger docs
- **Docker Support**: Ready-to-deploy containers
- **Smart Module Initialization**: Automatically initializes Go modules
- **Module Reference Replacement**: Updates all internal module references
- **Automatic Dependency Management**: Runs `go mod tidy` to download all packages

## 🎯 What You Get

Your generated project includes:

- ✅ Clean Architecture with dependency injection
- ✅ User authentication and authorization
- ✅ Real-time communication (SSE/WebSocket)
- ✅ Database models with GORM
- ✅ Email service with templates
- ✅ File upload and storage
- ✅ Admin dashboard endpoints
- ✅ Swagger API documentation
- ✅ Docker configuration
- ✅ Comprehensive README and setup instructions

## 📋 Requirements

- Go 1.19 or later
- Git

## 🛠️ Configuration Options

During setup, you can configure:

- **Project Name**: Your API project name
- **Go Module**: Module path (e.g., `github.com/username/project`)
- **Database**: SQLite (default), MySQL, or PostgreSQL
- **Port**: Server port (default: 8080)

## 🔧 Smart Module Management

GoInit includes intelligent module management features:

### Automatic Go Module Initialization

- **Checks for existing go.mod**: If a `go.mod` file already exists, it skips initialization
- **Runs `go mod init`**: Automatically initializes a new Go module with your specified module name
- **No manual setup required**: The generated project is ready to use immediately

### Module Reference Replacement

- **Scans all Go files**: Searches through all generated Go files for module references
- **Replaces old references**: Updates any hardcoded module references with your custom module name
- **Maintains consistency**: Ensures all internal imports use the correct module path

### Example

```bash
# User specifies: github.com/mycompany/myapi
# GoInit will:
# 1. Run: go mod init github.com/mycompany/myapi
# 2. Replace any "sog.com/goinit" references with "github.com/mycompany/myapi"
# 3. Run: go mod tidy (downloads all dependencies)
# 4. Generate a fully functional project ready to run
```

## 📁 Generated Project Structure

```
your-project/
├── cmd/api/                 # Application entry point
├── config/                  # Configuration management
├── internal/
│   ├── app/                # Application services
│   ├── data/               # Data layer (repositories)
│   ├── domain/             # Domain models
│   ├── lib/                # Shared libraries
│   └── server/             # Server setup
├── api/                    # HTTP handlers and routes
├── docs/                   # API documentation
├── .env                    # Environment configuration
├── go.mod                  # Go module file
├── README.md               # Project documentation
└── Dockerfile              # Docker configuration
```

## 🚀 Quick Start with Generated Project

```bash
# Navigate to your new project
cd your-project-name

# Dependencies are already downloaded by GoInit
# Configure environment (optional)
# Edit .env file with your settings

# Run the server
go run cmd/api/main.go
```

Visit `http://localhost:8080` and `http://localhost:8080/docs` for API documentation.

## 📚 API Endpoints

### Authentication

- `POST /api/auth/register/` - User registration
- `POST /api/auth/login/` - User login
- `GET /api/auth/logout/` - User logout

### Real-time

- `GET /api/sse/events` - Server-Sent Events
- `GET /api/ws/connect` - WebSocket connection

### Admin

- `GET /api/admin/users/` - List users
- `GET /api/admin/stats/` - User statistics

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- Inspired by Django's admin and authentication patterns

---

**Happy coding! 🎉**
# Test change for hook
# Test commit
