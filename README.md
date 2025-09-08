# Go Gin API Starter Template & Generator

A comprehensive Go API starter template built with Gin framework, featuring authentication, user management, real-time communication, and a CLI tool to generate new projects.

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

### 🏗️ Architecture

- **Clean Architecture** with separation of concerns
- **Dependency Injection** container
- **Repository Pattern** for data access
- **Service Layer** for business logic
- **Middleware** for cross-cutting concerns

## 🚀 Quick Start

### Option 1: Use the Template Directly

1. **Clone and setup:**

   ```bash
   git clone <repository-url>
   cd go-gin-api-template/gin
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Configure environment:**

   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

4. **Run the server:**
   ```bash
   go run cmd/api/main.go
   ```

### Option 2: Generate New Project with CLI

1. **Install the generator:**

   ```bash
   ./install.sh
   ```

2. **Generate new project:**

   ```bash
   goinit-generator
   ```

3. **Follow the prompts** to configure your project

4. **Start developing:**
   ```bash
   cd your-project-name
   go mod tidy
   go run cmd/api/main.go
   ```

## 📁 Project Structure

```
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
└── cli-generator/          # Project generator CLI
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

## 🐳 Docker Support

### Development

```bash
docker-compose up --build
```

### Production

```bash
# Build production image
docker build -f Dockerfile -t my-gin-api .

# Run with environment
docker run -p 8080:8080 --env-file .env my-gin-api
```

## 📚 API Documentation

When running, visit: http://localhost:8080/docs/

The documentation is auto-generated from code comments using Swagger.

## 🔧 Development

### Adding New Features

1. **Create domain models** in `internal/domain/`
2. **Implement repository** in `internal/data/`
3. **Add service logic** in `internal/app/`
4. **Create HTTP handler** in `api/protocol/http/handler/`
5. **Add routes** in `api/protocol/http/routes/`
6. **Register dependencies** in `internal/di/container.go`

### Code Generation

```bash
# Generate Swagger docs
swag init -g cmd/api/main.go

# Format code
go fmt ./...

# Vet code
go vet ./...
```

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

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [JWT](https://github.com/golang-jwt/jwt)
- [Swagger](https://swagger.io/)

---

**Happy coding! 🎉**

For questions or issues, please open a GitHub issue.
