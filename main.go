package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const version = "v1.0.0"

type ProjectConfig struct {
	ProjectName    string
	ModuleName     string
	DatabaseDriver string
	Port           string
}

func main() {
	// Check for version flag
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("GoInit %s\n", version)
		return
	}

	// Check for help flag
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		printHelp()
		return
	}

	fmt.Println("🚀 GoInit - Go Gin API Generator")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("=================================")

	// Get project configuration
	config := getProjectConfig()

	// Create project directory
	projectPath := config.ProjectName
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		fmt.Printf("❌ Error creating project directory: %v\n", err)
		return
	}

	fmt.Printf("📁 Creating project in: %s\n", projectPath)

	// Copy template files
	templatePath := "gin" // Path to the template
	if err := copyTemplate(templatePath, projectPath, config); err != nil {
		fmt.Printf("❌ Error copying template: %v\n", err)
		return
	}

	// Generate project files with templating
	if err := generateTemplatedFiles(projectPath, config); err != nil {
		fmt.Printf("❌ Error generating templated files: %v\n", err)
		return
	}

	// Initialize Go module if not present
	if err := initializeGoModule(projectPath, config); err != nil {
		fmt.Printf("❌ Error initializing Go module: %v\n", err)
		return
	}

	// Replace module references in all Go files
	if err := replaceModuleReferences(projectPath, config); err != nil {
		fmt.Printf("❌ Error replacing module references: %v\n", err)
		return
	}

	// Run go mod tidy to download dependencies
	if err := runGoModTidy(projectPath); err != nil {
		fmt.Printf("❌ Error running go mod tidy: %v\n", err)
		return
	}

	fmt.Println("\n✅ Project generated successfully!")
	fmt.Printf("📁 Project location: %s\n", projectPath)
	fmt.Println("\n🚀 Next steps:")
	fmt.Printf("  cd %s\n", projectPath)
	fmt.Println("  go run cmd/api/main.go")
	fmt.Println("\n📚 API Documentation: http://localhost:8080/docs")
	fmt.Println("🎉 Happy coding!")
}

func printHelp() {
	fmt.Println("GoInit - Go Gin API Generator")
	fmt.Printf("Version: %s\n", version)
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  goinit [flags]")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  --version, -v    Show version information")
	fmt.Println("  --help, -h       Show this help message")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  Interactive CLI tool to generate production-ready Go API projects")
	fmt.Println("  built with the Gin framework, featuring authentication, real-time")
	fmt.Println("  communication, and comprehensive API documentation.")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  goinit              # Generate a new project interactively")
	fmt.Println("  goinit --version    # Show version")
	fmt.Println("  goinit --help       # Show this help")
}

func getProjectConfig() ProjectConfig {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter project name: ")
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	fmt.Print("Enter Go module name (e.g., github.com/username/project): ")
	moduleName, _ := reader.ReadString('\n')
	moduleName = strings.TrimSpace(moduleName)

	if moduleName == "" {
		moduleName = fmt.Sprintf("github.com/user/%s", projectName)
	}

	fmt.Print("Choose database driver (sqlite/mysql/postgres) [sqlite]: ")
	dbDriver, _ := reader.ReadString('\n')
	dbDriver = strings.TrimSpace(dbDriver)
	if dbDriver == "" {
		dbDriver = "sqlite"
	}

	fmt.Print("Enter port [8080]: ")
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)
	if port == "" {
		port = "8080"
	}

	return ProjectConfig{
		ProjectName:    projectName,
		ModuleName:     moduleName,
		DatabaseDriver: dbDriver,
		Port:           port,
	}
}

func copyTemplate(src, dst string, config ProjectConfig) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip certain directories and files
		if shouldSkip(path, info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func shouldSkip(path string, info os.FileInfo) bool {
	// Skip hidden files and directories
	if strings.HasPrefix(info.Name(), ".") {
		return true
	}

	// Skip build artifacts
	if info.Name() == "tmp" || info.Name() == "logs" {
		return true
	}

	// Skip specific files
	skipFiles := []string{
		"go.sum",
		".env",
		"README.md", // We'll generate our own
	}

	for _, skip := range skipFiles {
		if info.Name() == skip {
			return true
		}
	}

	return false
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func generateTemplatedFiles(projectPath string, config ProjectConfig) error {
	// Note: go.mod will be created by initializeGoModule function

	// Generate .env file
	if err := generateEnvFile(projectPath, config); err != nil {
		return err
	}

	// Generate README
	if err := generateReadme(projectPath, config); err != nil {
		return err
	}

	return nil
}

func generateEnvFile(projectPath string, config ProjectConfig) error {
	envContent := fmt.Sprintf(`# Server Configuration
PORT=%s
PUBLIC_HOST=http://localhost:%s

# Database Configuration
DB_DRIVER=%s
DB_USER=root
DB_PASSWORD=password
DB_NAME=%s
DB_HOST=127.0.0.1
DB_PORT=3306

# Session Configuration
SESSION_SECRET=dev-session-secret-change-me-in-production
SESSION_NAME=hor_session
SESSION_SECURE=false
SESSION_DOMAIN=
SESSION_MAX_AGE=86400

# JWT Configuration
JWT_SECRET=dev-jwt-secret-change-me-in-production
USE_DATABASE_JWT=false

# Email Configuration
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=
EMAIL_PASSWORD=
EMAIL_FROM=noreply@%s.com
USE_LOCAL_EMAIL=true
EMAIL_LOG_PATH=./logs/emails.log

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Password Reset Configuration
USE_DATABASE_PWRESET=false

# Storage Configuration
STORAGE_BACKEND=local
UPLOAD_BASE_DIR=./uploads
UPLOAD_PUBLIC_BASE_URL=/uploads

# S3 Configuration (if using S3 storage)
S3_ENDPOINT=
S3_REGION=us-east-1
S3_BUCKET=
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY=
S3_USE_SSL=true
S3_FORCE_PATH_STYLE=false
S3_PUBLIC_BASE_URL=

# Logging
LOG_LEVEL=info
LOG_FILE=logs/app.log
LOG_FILE_ENABLED=false
GIN_MODE=debug
`, config.Port, config.Port, config.DatabaseDriver, config.ProjectName, config.ProjectName)

	return os.WriteFile(filepath.Join(projectPath, ".env"), []byte(envContent), 0644)
}

func generateReadme(projectPath string, config ProjectConfig) error {
	readmeContent := fmt.Sprintf(`# %s

A Go API server built with Gin framework, featuring authentication, user management, real-time communication, and more.

## Features

- 🔐 **Authentication & Authorization**
  - JWT-based authentication
  - Session management
  - Password reset functionality
  - Admin role management

- 👥 **User Management**
  - User registration and login
  - Profile management
  - User roles and permissions
  - Admin user controls

- 📡 **Real-time Communication**
  - Server-Sent Events (SSE)
  - WebSocket support
  - Real-time notifications

- 🗄️ **Database Support**
  - SQLite (default)
  - MySQL
  - PostgreSQL

- 📧 **Email Integration**
  - SMTP email sending
  - Local email logging for development
  - Password reset emails

- ☁️ **Storage**
  - Local file storage
  - S3-compatible storage

- 📚 **API Documentation**
  - Swagger/OpenAPI documentation
  - Auto-generated docs

## Quick Start

1. **Install dependencies:**
   ` + "```" + `bash
   go mod tidy
   ` + "```" + `

2. **Set up environment:**
   ` + "```" + `bash
   cp .env.example .env
   # Edit .env with your configuration
   ` + "```" + `

3. **Run the server:**
   ` + "```" + `bash
   go run cmd/api/main.go
   ` + "```" + `

The server will start on http://localhost:%s

## API Endpoints

### Authentication
- ` + "`POST /api/auth/register/`" + ` - User registration
- ` + "`POST /api/auth/login/`" + ` - User login
- ` + "`POST /api/auth/logout/`" + ` - User logout
- ` + "`POST /api/auth/change-password/`" + ` - Change password

### User Management
- ` + "`GET /api/user/profile/`" + ` - Get user profile
- ` + "`PUT /api/user/profile/`" + ` - Update user profile

### Real-time
- ` + "`GET /api/sse/events`" + ` - Server-Sent Events stream
- ` + "`GET /api/ws/connect`" + ` - WebSocket connection

### Admin
- ` + "`GET /api/admin/users/`" + ` - List all users
- ` + "`GET /api/admin/stats/`" + ` - User statistics

## Configuration

The application uses environment variables for configuration. Copy .env.example to .env and modify as needed.

Key configuration options:
- ` + "`DB_DRIVER`" + `: Database driver (sqlite/mysql/postgres)
- ` + "`JWT_SECRET`" + `: JWT signing secret
- ` + "`SESSION_SECRET`" + `: Session signing secret
- ` + "`EMAIL_*`" + `: Email configuration
- ` + "`REDIS_*`" + `: Redis configuration

## Development

### Project Structure
` + "```" + `
├── cmd/api/           # Application entry point
├── config/            # Configuration management
├── internal/
│   ├── app/          # Application logic
│   ├── data/         # Data layer (repositories)
│   ├── domain/       # Domain models
│   ├── lib/          # Shared libraries
│   └── server/       # Server setup
├── api/              # API handlers and routes
└── docs/             # API documentation
` + "```" + `

### Building

` + "```" + `bash
go build -o bin/server cmd/api/main.go
` + "```" + `

### Testing

` + "```" + `bash
go test ./...
` + "```" + `

## Docker Support

Build and run with Docker:

` + "```" + `bash
docker-compose up --build
` + "```" + `

## API Documentation

Once the server is running, visit http://localhost:%s/docs/ for Swagger documentation.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.
`, config.ProjectName, config.Port, config.Port)

	return os.WriteFile(filepath.Join(projectPath, "README.md"), []byte(readmeContent), 0644)
}

// initializeGoModule checks if go.mod exists and runs go mod init if not
func initializeGoModule(projectPath string, config ProjectConfig) error {
	goModPath := filepath.Join(projectPath, "go.mod")

	// Check if go.mod already exists
	if _, err := os.Stat(goModPath); err == nil {
		fmt.Println("📦 go.mod already exists, skipping initialization")
		return nil
	}

	fmt.Printf("📦 Initializing Go module: %s\n", config.ModuleName)

	// Run go mod init
	cmd := exec.Command("go", "mod", "init", config.ModuleName)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize Go module: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("✅ Go module initialized successfully")
	return nil
}

// replaceModuleReferences replaces all module references in Go files
func replaceModuleReferences(projectPath string, config ProjectConfig) error {
	fmt.Println("🔄 Replacing module references in Go files...")

	// Find all Go files in the project
	var goFiles []string
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk project directory: %v", err)
	}

	// Replace module references in each Go file
	replacements := 0
	for _, filePath := range goFiles {
		if replaced, err := replaceInFile(filePath, "github.com/SOG-web/goinit/gin", config.ModuleName); err != nil {
			return fmt.Errorf("failed to replace in file %s: %v", filePath, err)
		} else if replaced {
			replacements++
		}
	}

	if replacements > 0 {
		fmt.Printf("✅ Replaced module references in %d files\n", replacements)
	} else {
		fmt.Println("ℹ️  No module references found to replace")
	}

	return nil
}

// replaceInFile replaces all occurrences of oldString with newString in the given file
func replaceInFile(filePath, oldString, newString string) (bool, error) {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	// Check if the file contains the old string
	if !strings.Contains(string(content), oldString) {
		return false, nil
	}

	// Replace all occurrences
	newContent := strings.ReplaceAll(string(content), oldString, newString)

	// Write the file back
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}

// runGoModTidy runs go mod tidy to download dependencies and update go.sum
func runGoModTidy(projectPath string) error {
	fmt.Println("📦 Running go mod tidy to download dependencies...")

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("✅ Dependencies downloaded successfully")
	return nil
}