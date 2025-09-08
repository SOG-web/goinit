# GoInit - Go Gin API Starter Template & Generator

A comprehensive Go API starter template built with Gin framework, featuring authentication, user management, real-time communication, and a CLI tool to generate new projects from the template.

## 🎯 What's Included

### ✅ Core Features

- **Authentication & Authorization**
  - JWT-based authentication with refresh tokens
  - Session management with cookies
  - Password reset functionality
  - Role-based access control (admin/user)

- **User Management**
  - User registration and login
  - Profile management
  - User roles and permissions
  - Admin user controls

- **Real-time Communication**
  - Server-Sent Events (SSE) for notifications
  - WebSocket support for bidirectional communication
  - Real-time event streaming

- **Database Support**
  - SQLite (default, file-based)
  - MySQL
  - PostgreSQL
  - GORM ORM with auto-migration

- **Email Integration**
  - SMTP email sending
  - Local email logging for development
  - HTML email templates

- **Storage Solutions**
  - Local file storage
  - S3-compatible storage (AWS S3, MinIO, etc.)

- **API Documentation**
  - Swagger/OpenAPI 3.0 documentation
  - Auto-generated API docs
  - Interactive API testing

- **Project Generator CLI**
  - Interactive project setup
  - Customizable configuration
  - Automated dependency management
  - Cross-platform builds

### 🏗️ Architecture

- **Clean Architecture** with separation of concerns
- **Dependency Injection** container
- **Repository Pattern** for data access
- **Service Layer** for business logic
- **Middleware** for cross-cutting concerns
- **Automated Git Hooks** for template synchronization

## 🚀 Quick Start

### Option 1: Use the Template Directly

1. **Clone the repository:**

   ```bash
   git clone https://github.com/SOG-web/goinit.git
   cd goinit
   ```

2. **Navigate to the template:**

   ```bash
   cd gin
   ```

3. **Install dependencies:**

   ```bash
   go mod tidy
   ```

4. **Configure environment:**

   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

5. **Run the server:**
   ```bash
   go run cmd/api/main.go
   ```

### Option 2: Generate New Project with CLI (Recommended)

1. **Clone the repository:**

   ```bash
   git clone https://github.com/SOG-web/goinit.git
   cd goinit
   ```

2. **Install the generator:**

   ```bash
   chmod +x install.sh
   ./install.sh
   ```

3. **Generate new project:**

   ```bash
   goinit-generator
   ```

4. **Follow the interactive prompts** to configure your project

5. **Start developing:**
   ```bash
   cd your-project-name
   go mod tidy
   go run cmd/api/main.go
   ```

### Option 3: Download Pre-built Binary

1. **Go to [Releases](https://github.com/SOG-web/goinit/releases)**
2. **Download the appropriate binary** for your platform
3. **Make it executable and run:**

   ```bash
   chmod +x goinit-*
   ./goinit-* --help
   ```

## 📁 Project Structure

### Repository Structure

This repository contains both the template and the CLI generator:

```
goinit/
├── gin/                     # 🏗️  API Template Source
│   ├── cmd/api/            # Application entry point
│   ├── config/             # Configuration management
│   ├── internal/           # Internal application code
│   ├── api/                # HTTP handlers and routes
│   ├── docs/               # API documentation
│   └── docker/             # Docker configuration
├── cli-generator/          # 🛠️  CLI Tool Source
│   ├── main.go            # CLI entry point
│   ├── .git/hooks/        # Git hooks for template sync
│   └── gin/               # Copied template for CLI use
├── .github/workflows/      # 🚀 GitHub Actions CI/CD
├── go.mod                 # Go module for CLI generator
├── install.sh             # Installation script
└── README.md              # This file
```

### Generated Project Structure

When you use the CLI generator, it creates a new project with this structure:

```
your-project/
├── cmd/api/                 # Application entry point
├── config/                  # Configuration management
├── internal/
│   ├── app/                # Application services
│   ├── data/               # Data layer (repositories)
│   ├── domain/             # Domain models
│   ├── lib/                # Shared libraries
│   │   ├── email/          # Email service
│   │   ├── jwt/            # JWT utilities
│   │   ├── pwreset/        # Password reset
│   │   └── storage/        # File storage
│   └── server/             # Server setup
├── api/                    # HTTP handlers and routes
│   ├── common/
│   │   ├── dto/            # Data transfer objects
│   │   └── middleware/     # HTTP middleware
│   └── protocol/
│       ├── http/
│       │   ├── handler/    # HTTP handlers
│       │   ├── router/     # Route setup
│       │   └── routes/     # Route definitions
│       ├── sse/            # Server-Sent Events
│       └── ws/             # WebSocket handlers
├── docs/                   # API documentation
├── docker/                 # Docker configuration
├── .env.example           # Environment template
└── go.mod                 # Go module file
```

## 🔧 Configuration

The application uses environment variables for configuration. Key settings:

### Server

```env
PORT=8080
PUBLIC_HOST=http://localhost:8080
GIN_MODE=debug
```

### Database

```env
DB_DRIVER=sqlite          # sqlite, mysql, postgres
DB_NAME=myapp
DB_USER=root
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=3306
```

### Authentication

```env
JWT_SECRET=your-jwt-secret
SESSION_SECRET=your-session-secret
USE_DATABASE_JWT=false    # true for database, false for Redis
```

### Email

```env
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
USE_LOCAL_EMAIL=true      # true for development logging
```

### Real-time Features

```env
REDIS_ADDR=localhost:6379  # For JWT blacklist and sessions
```

## 📡 API Endpoints

### Authentication

- `POST /api/auth/register/` - User registration
- `POST /api/auth/login/` - User login
- `GET /api/auth/logout/` - User logout
- `POST /api/auth/change-password/` - Change password
- `POST /api/auth/password-reset/request/` - Request password reset
- `POST /api/auth/password-reset/confirm/` - Confirm password reset

### User Management

- `GET /api/user/profile/` - Get user profile
- `PUT /api/user/profile/` - Update user profile
- `POST /api/user/profile/image/` - Upload profile image

### Real-time

- `GET /api/sse/events` - Server-Sent Events stream
- `GET /api/sse/notifications` - Notification SSE stream
- `GET /api/ws/connect` - WebSocket connection

### Admin (requires admin role)

- `GET /api/admin/users/` - List all users
- `GET /api/admin/stats/` - User statistics
- `PUT /api/admin/users/:id/activate/` - Activate user
- `PUT /api/admin/users/:id/deactivate/` - Deactivate user

### Health Check

- `GET /health` - Server health check

## 🔐 Authentication

The API supports multiple authentication methods:

### JWT Bearer Token

```bash
curl -H "Authorization: Bearer <jwt-token>" \
     http://localhost:8080/api/user/profile/
```

### Session Cookies

Session cookies are automatically managed by the session middleware.

## 📧 Email Templates

The application includes email templates for:

- User registration verification
- Password reset
- Welcome emails
- Admin notifications

## 🗄️ Database Migrations

The application uses GORM's auto-migration feature. Models are automatically migrated on startup.

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestUserService ./internal/app/user/
```

## � GitHub Actions & Releases

This project includes automated CI/CD pipelines:

### Automated Releases

- **Trigger**: Push a version tag (e.g., `v1.0.0`)
- **Builds**: Cross-platform binaries for Linux, macOS (Intel/Apple Silicon), and Windows
- **Artifacts**: Automatically uploaded to GitHub Releases
- **Workflow**: `.github/workflows/release.yml`

### Creating a Release

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will automatically:
# 1. Build binaries for all platforms
# 2. Create release archives
# 3. Upload to GitHub Releases
```

### Development Workflow

- **Pre-commit hooks**: Automatically sync template files
- **Pre-push hooks**: Ensure template consistency before pushing
- **Automated testing**: Run tests on all platforms
- **Code quality**: Automated linting and formatting

## � Docker Support

### Development

```bash
# From the template directory
cd gin
docker-compose up --build
```

### Production

```bash
# From the template directory
cd gin

# Build production image
docker build -f Dockerfile -t my-gin-api .

# Run with environment
docker run -p 8080:8080 --env-file .env my-gin-api
```

## �📚 API Documentation

When running the generated project, visit: http://localhost:8080/docs/

The documentation is auto-generated from code comments using Swagger.

## 🔧 Development

### Adding New Features

1. **Create domain models** in `internal/domain/`
2. **Implement repository** in `internal/data/`
3. **Add service logic** in `internal/app/`
4. **Create HTTP handler** in `api/protocol/http/handler/`
5. **Add routes** in `api/protocol/http/routes/`
6. **Register dependencies** in `internal/di/container.go`

### Git Hooks

This project includes automated Git hooks for template synchronization:

- **Pre-commit**: Automatically copies updated template files before commits
- **Pre-push**: Ensures template consistency before pushing changes
- **Location**: `.git/hooks/` (automatically installed)

The hooks ensure that the CLI generator always uses the latest template files.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Update documentation
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [JWT](https://github.com/golang-jwt/jwt) - JSON Web Tokens
- [Swagger](https://swagger.io/) - API documentation
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Gorilla WebSocket](https://github.com/gorilla/websocket) - WebSocket implementation
- [GoMail](https://github.com/go-gomail/gomail) - Email sending

---

**Repository**: [https://github.com/SOG-web/goinit](https://github.com/SOG-web/goinit)  
**Happy coding! 🎉**

For questions or issues, please open a GitHub issue.
