# PostgreSQL-Golang API Project

A modern REST API built with Go, PostgreSQL, and Gin framework with comprehensive authentication, role-based access control, and AWS S3 integration.

## 📋 Overview

This project is a production-ready backend API that provides:
- User authentication and authorization with JWT tokens
- Role-based access control (RBAC)
- Admin management system
- Staff management
- AWS S3 file storage integration
- Request rate limiting and security middleware
- Database migrations and seeding

## 🚀 Tech Stack

- **Language**: Go 1.26.3
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL with [GORM](https://gorm.io/) ORM
- **Authentication**: JWT (JSON Web Tokens)
- **Cloud Storage**: AWS S3
- **Database Migrations**: Custom migration system
- **Logging**: Go's structured logging (slog)

## 📦 Key Dependencies

- `gin-gonic/gin` - HTTP web framework
- `gorm.io/gorm` - ORM for database operations
- `golang-jwt/jwt` - JWT token generation and validation
- `aws-sdk-go-v2` - AWS S3 integration
- `gin-contrib/cors` - CORS middleware
- `gin-contrib/rate-limiting` - Rate limiting
- `gin-contrib/secure` - Security headers
- `golang.org/x/crypto` - Cryptographic functions

## 📁 Project Structure

```
.
├── bin/                          # Compiled binaries
│   └── api                       # API executable
├── cmd/                          # Command-line applications
│   ├── api/                      # Main API server
│   ├── keygen/                   # Key generation utility
│   └── seeder/                   # Database seeder
├── internal/                     # Internal packages
│   ├── bootstrap/                # Dependency injection container
│   ├── common/                   # Common utilities
│   ├── config/                   # Configuration management
│   │   ├── app.go               # App configuration
│   │   ├── database.go          # Database configuration
│   │   └── s3.go                # AWS S3 configuration
│   ├── constants/                # Application constants
│   ├── database/                 # Database layer
│   │   ├── migrations/           # Database migrations
│   │   └── seeders/              # Database seeders
│   ├── helpers/                  # Helper functions
│   │   ├── helper.go
│   │   ├── jwt.go               # JWT utilities
│   │   └── s3-helper.go         # S3 utilities
│   ├── logger/                   # Logging configuration
│   ├── middlewares/              # HTTP middleware
│   │   ├── auth.go              # Authentication
│   │   ├── permission.go        # Permission checking
│   │   ├── rate-limiting.go     # Rate limiting
│   │   ├── size.go              # Request size limiting
│   │   ├── throttle.go          # Throttling
│   │   └── timeout.go           # Request timeout
│   ├── modules/                  # Feature modules
│   │   ├── auth/                 # Authentication module
│   │   ├── health/               # Health check
│   │   ├── permission/           # Permission model
│   │   ├── personalaccesstoken/  # Token management
│   │   ├── role/                 # Role management
│   │   └── staff/                # Staff management
│   └── routes/                   # Route definitions
├── go.mod                        # Go module dependencies
├── Makefile                      # Build automation
└── README.md                     # This file
```

## 🔧 Installation

### Prerequisites

- Go 1.26.3 or higher
- PostgreSQL 12 or higher
- AWS account (for S3 integration, optional)

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Postgresql-golang-project
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your configuration:
   ```env
   APP_NAME=Your App Name
   APP_ENV=local
   APP_DEBUG=true
   APP_KEY=your-secret-key
   APP_URL=http://localhost
   APP_PORT=8080
   
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=password
   DB_NAME=your_database
   DB_SSL_MODE=disable
   
   # JWT
   JWT_ACCESS_SECRET=your-jwt-access-secret
   JWT_REFRESH_SECRET=your-jwt-refresh-secret
   JWT_ACCESS_EXPIRY_MINUTES=60
   JWT_REFRESH_EXPIRY_DAYS=30
   
   # AWS S3
   AWS_REGION=us-east-1
   AWS_ACCESS_KEY_ID=your-access-key
   AWS_SECRET_ACCESS_KEY=your-secret-key
   S3_BUCKET=your-bucket-name
   
   # CORS
   ALLOWED_ORIGIN=http://localhost:3000
   ALLOWED_HOSTS=localhost
   ```

4. **Create the database**
   ```bash
   createdb your_database
   ```

5. **Run migrations**
   ```bash
   make seed
   ```

## 🏃 Running the Application

### Development

```bash
# Run with automatic reload
make dev

# Or run the API directly
go run ./cmd/api
```

### Production

```bash
# Build binary
make build

# Run the binary
make start
```

### Other Commands

```bash
# Run database seeder
make seed

# Generate API keys
make keygen

# Format code
make fmt

# Run static analysis
make vet

# Run tests
make test

# Clean build artifacts
make clean

# View all available commands
make help
```

## 🔐 Authentication & Authorization

### JWT Tokens

The API uses JWT for authentication with two types of tokens:
- **Access Token**: Short-lived token (default 60 minutes)
- **Refresh Token**: Long-lived token (default 30 days)

### Role-Based Access Control

The system implements RBAC with different permission levels:
- **Admin**: Full system access
- **Staff**: Limited access based on assigned permissions
- **User**: Basic access

### Middleware

- **AuthMiddleware**: Validates JWT tokens
- **GuardMiddleware**: Checks user roles and permissions
- **PermissionMiddleware**: Fine-grained permission checking
- **RateLimitMiddleware**: Prevents abuse with rate limiting
- **TimeoutMiddleware**: Sets request timeouts
- **SizeMiddleware**: Limits request body size
- **ThrottleMiddleware**: Throttles requests
- **CORSMiddleware**: Handles cross-origin requests
- **SecureMiddleware**: Adds security headers

## 📚 API Modules

### Auth Module
- User login/logout
- Token refresh
- Password management
- Admin authentication

### Role Module
- CRUD operations for roles
- Permission assignment
- Role hierarchy management

### Staff Module
- Staff member management
- Profile updates
- Activity tracking

### Health Check
- API health status endpoint

## 🗄️ Database

### Models

- **User**: User accounts and authentication
- **Role**: Role definitions and hierarchy
- **Permission**: Permission registry
- **Staff**: Staff member information
- **PersonalAccessToken**: API tokens for automation

### Migrations

Database migrations are automatically run on startup. Custom migration files are located in `internal/database/migrations/`.

### Seeders

Seed the database with initial data:
```bash
make seed
```

Available seeders:
- Admin seeder: Creates default admin user and roles
- Permission seeder: Populates permission registry

## 📦 AWS S3 Integration

The project includes S3 integration for file storage:

```go
// Upload file
fileURL, err := s3Helper.UploadFile(ctx, file)

// Delete file
err := s3Helper.DeleteFile(ctx, fileKey)
```

Configure S3 in `.env`:
```env
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-key
AWS_SECRET_ACCESS_KEY=your-secret
S3_BUCKET=your-bucket
```

## 🧪 Testing

Run all tests:
```bash
make test
```

Tests are located throughout the project alongside the code they test.

## 📝 Logging

The project uses Go's structured logging (slog). Logs are configured in `internal/logger/logger.go`.

Log levels:
- **Info**: General information
- **Warn**: Warning messages
- **Error**: Error messages
- **Debug**: Debug information (development only)

## 🛡️ Security Features

- JWT-based authentication
- Password hashing with bcrypt
- CORS policy enforcement
- Rate limiting
- Request size limiting
- Security headers (HTTPS, HSTS, etc.)
- SQL injection prevention (via GORM)
- CSRF protection
- Request timeout enforcement

## 🤝 Contributing

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make your changes and commit: `git commit -am 'Add your feature'`
3. Push to the branch: `git push origin feature/your-feature`
4. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 👤 Author

Created by Akshit Tyagi

## 🆘 Support

For issues, questions, or suggestions, please open an issue in the repository.

---

**Last Updated**: 2026-08-20
