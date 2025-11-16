# Multi-Tenant Chat Application

A production-ready, scalable multi-tenant chat application built with Go, PostgreSQL, and Stream Chat API. Features role-based access control (RBAC), JWT authentication, and automated AWS Fargate deployment.

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![AWS](https://img.shields.io/badge/AWS-Fargate-orange.svg)](https://aws.amazon.com/fargate/)

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Deployment](#deployment)
- [Monitoring](#monitoring)
- [Contributing](#contributing)

## ✨ Features

### Core Functionality
- **Multi-Tenancy**: Isolated data per organization with tenant-based access control
- **Real-time Messaging**: Powered by Stream Chat API for instant communication
- **Role-Based Access Control (RBAC)**: Four roles - Admin, Moderator, Member, Guest
- **JWT Authentication**: Secure token-based authentication
- **Channel Management**: Create, join, and manage chat channels
- **User Management**: Full CRUD operations with role-based permissions

### Technical Features
- **RESTful API**: Clean, well-documented API endpoints
- **Database Mocking**: Comprehensive test suite with sqlmock
- **Infrastructure as Code**: Terraform modules for AWS deployment
- **CI/CD Pipeline**: Automated testing and deployment via GitHub Actions
- **Health Checks**: Built-in health monitoring endpoints
- **CORS Support**: Configured for frontend integration
- **Swagger Documentation**: Interactive API documentation

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Internet                             │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │  Application Load    │
              │     Balancer (ALB)   │
              └──────────┬───────────┘
                         │
         ┌───────────────┴───────────────┐
         │                               │
         ▼                               ▼
┌─────────────────┐           ┌─────────────────┐
│  ECS Fargate    │           │  ECS Fargate    │
│   Task (AZ-1)   │           │   Task (AZ-2)   │
│   Go Backend    │           │   Go Backend    │
└────────┬────────┘           └────────┬────────┘
         │                               │
         └───────────────┬───────────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │   RDS PostgreSQL     │
              │   (Multi-AZ Ready)   │
              └──────────────────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │   Stream Chat API    │
              │   (External SaaS)    │
              └──────────────────────┘
```

### Key Components

- **Backend**: Go 1.24 with Gin framework
- **Database**: PostgreSQL 15.7 (AWS RDS)
- **Real-time**: Stream Chat API
- **Container Orchestration**: AWS ECS Fargate
- **Load Balancing**: Application Load Balancer
- **Secrets Management**: AWS Secrets Manager
- **Monitoring**: CloudWatch Logs & Metrics

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.24
- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **Real-time**: Stream Chat Go SDK

### Database
- **Primary**: PostgreSQL 15.7
- **Driver**: pgx/v5

### Infrastructure
- **Cloud Provider**: AWS
- **Compute**: ECS Fargate (serverless containers)
- **Database**: RDS PostgreSQL
- **Load Balancer**: Application Load Balancer
- **Container Registry**: ECR
- **IaC**: Terraform
- **CI/CD**: GitHub Actions

### Testing
- **Framework**: testify
- **Mocking**: sqlmock (database mocking)

## 📁 Project Structure

```
chat-app/
├── cmd/
│   └── main.go                 # Application entry point
├── db/
│   ├── db.go                   # Database connection
│   └── mock.go                 # Database mocking utilities
├── handlers/
│   ├── auth.go                 # Authentication handlers
│   ├── channel.go              # Channel management
│   ├── channel_test.go         # Channel tests
│   ├── stream.go               # Stream Chat integration
│   └── tenant_user.go          # Tenant & user management
├── middleware/
│   ├── auth.go                 # Authentication middleware
│   ├── jwt.go                  # JWT utilities
│   └── rbac.go                 # Role-based access control
├── models/
│   └── models.go               # Data models
├── services/
│   ├── bcrypt.go               # Password hashing
│   ├── channel.go              # Channel business logic
│   └── stream.go               # Stream Chat service
├── utils/
│   └── jwt.go                  # JWT helper functions
├── testutil/
│   └── helpers.go              # Test utilities
├── terraform/
│   ├── main.tf                 # Root Terraform config
│   ├── variables.tf            # Input variables
│   ├── outputs.tf              # Output values
│   └── modules/
│       ├── vpc/                # Network infrastructure
│       ├── rds/                # Database module
│       ├── ecs/                # Container orchestration
│       ├── alb/                # Load balancer
│       └── secrets/            # Secrets management
├── .github/
│   └── workflows/
│       └── deploy.yml          # CI/CD pipeline
├── Dockerfile                  # Container image definition
├── docker-compose.yml          # Local development setup
├── setup-infrastructure.sh     # Infrastructure setup script
├── .env.example                # Environment variables template
└── README.md                   # This file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL 15+ (or use Docker Compose)
- Stream Chat account ([getstream.io](https://getstream.io))
- AWS CLI (for deployment)
- Terraform 1.0+ (for deployment)

### Local Development Setup

1. **Clone the repository**
```bash
git clone https://github.com/Nyagar-Abraham/chat-app.git
cd chat-app
```

2. **Set up environment variables**
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```env
DATABASE_URL=postgresql://chatuser:chatpassword@localhost:5434/chat_db?sslmode=disable
STREAM_API_KEY=your_stream_api_key
STREAM_API_SECRET=your_stream_api_secret
JWT_SECRET=your_64_character_secret
MIGRATE_DB=true
PORT=8085
```

3. **Start PostgreSQL with Docker Compose**
```bash
docker-compose up -d
```

4. **Install dependencies**
```bash
go mod download
```

5. **Run database migrations**
```bash
# Migrations run automatically when MIGRATE_DB=true
go run cmd/main.go
```

6. **Start the server**
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8085`

### Generate JWT Secret

```bash
openssl rand -hex 32
```

## 📚 API Documentation

### Interactive Documentation

Access Swagger UI at: `http://localhost:8085/swagger/index.html`

### Authentication

All authenticated endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### Core Endpoints

#### Authentication
```http
POST   /auth/register          # Register new user
POST   /auth/login             # Login user
GET    /me                     # Get current user
```

#### Tenants
```http
POST   /tenants                # Create tenant (Admin only)
GET    /tenants                # List all tenants
GET    /tenants/:id            # Get tenant by ID
```

#### Users
```http
POST   /users                  # Create user (Admin/Moderator)
GET    /users                  # List users
PUT    /users/:id              # Update user (Admin/Moderator)
DELETE /users/:id              # Delete user (Admin only)
```

#### Channels
```http
POST   /channels               # Create channel (Admin/Moderator)
GET    /channels               # List channels
POST   /channels/:id/join      # Join channel
POST   /channels/:id/leave     # Leave channel
GET    /channels/:id/members   # Get channel members
POST   /channels/:id/members   # Add user to channel (Admin/Moderator)
DELETE /channels/:id/members/:user_id  # Remove user (Admin/Moderator)
```

#### Messages
```http
POST   /messages               # Send message
GET    /messages/:stream_id    # Get messages
```

#### Stream Chat
```http
GET    /stream/token           # Get Stream Chat token
```

#### Health Check
```http
GET    /health                 # Health check endpoint
```

### Example Requests

**Register User:**
```bash
curl -X POST http://localhost:8085/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword",
    "role": "MEMBER",
    "org_name": "Acme Corp"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8085/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

**Create Channel:**
```bash
curl -X POST http://localhost:8085/channels \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "general",
    "description": "General discussion"
  }'
```

## 🧪 Testing

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Package Tests

```bash
go test ./handlers -v
```

### Testing Features

- **Database Mocking**: All tests use sqlmock for isolated testing
- **No Real Database Required**: Tests run without PostgreSQL
- **Fast Execution**: No I/O overhead
- **Comprehensive Coverage**: Handler, service, and middleware tests

### Test Structure

```go
func TestExample(t *testing.T) {
    mock := testutil.SetupMockDB(t)
    
    // Define expectations
    mock.ExpectQuery("SELECT...").WillReturnRows(...)
    
    // Run test
    // ...
    
    // Verify
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

See [TESTING_GUIDE.md](TESTING_GUIDE.md) for detailed testing documentation.

## 🚀 Deployment

### Deployment Architecture

The application uses a **two-phase deployment approach**:

1. **Infrastructure Setup** (One-time) - Terraform creates AWS resources
2. **Application Deployment** (Automated) - GitHub Actions deploys code

### Phase 1: Infrastructure Setup

#### Prerequisites

- AWS account with appropriate permissions
- AWS CLI configured
- Terraform installed
- Domain name (optional)

#### Step 1: Configure Terraform Variables

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars`:
```hcl
aws_region  = "us-east-1"
app_name    = "chat-app"
environment = "prod"

db_password        = "your_strong_password"
jwt_secret         = "your_64_char_secret"
stream_api_key     = "your_stream_key"
stream_api_secret  = "your_stream_secret"
```

#### Step 2: Run Infrastructure Setup

```bash
./setup-infrastructure.sh
```

This creates:
- VPC with public/private subnets (2 AZs)
- RDS PostgreSQL database
- ECS Fargate cluster
- Application Load Balancer
- ECR repository
- Secrets Manager secrets
- CloudWatch logs
- IAM roles and policies

#### Step 3: Update Database Secret

```bash
# Get database endpoint
DB_ENDPOINT=$(cd terraform && terraform output -raw db_endpoint)

# Update secret
aws secretsmanager update-secret \
  --secret-id chat-app/prod/database-url \
  --secret-string "postgresql://chatadmin:YOUR_PASSWORD@${DB_ENDPOINT}/chat_db?sslmode=require"
```

### Phase 2: CI/CD Deployment

#### Configure GitHub Secrets

Add these secrets to your GitHub repository:

```
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
```

#### Automatic Deployment

Push to `main` branch triggers automatic deployment:

```bash
git add .
git commit -m "Deploy to production"
git push origin main
```

#### Deployment Pipeline

```
Code Push → Tests → Build Docker → Push to ECR → Update ECS → Live
```

The GitHub Actions workflow:
1. Runs all tests with mocked database
2. Builds Docker image
3. Pushes to Amazon ECR
4. Updates ECS service
5. Waits for deployment to stabilize

### Infrastructure Details

**AWS Resources Created:**
- **VPC**: 10.0.0.0/16 with 2 AZs
- **Subnets**: 2 public, 2 private
- **NAT Gateways**: 2 (high availability)
- **RDS**: PostgreSQL 15.7, db.t3.micro
- **ECS**: Fargate cluster with 2 tasks
- **ALB**: Application Load Balancer
- **ECR**: Docker image registry

**Monthly Cost Estimate:** ~$130
- ECS Fargate: $15
- RDS: $15
- ALB: $20
- NAT Gateways: $65
- Other: $15

### Deployment Documentation

- [AWS_FARGATE_DEPLOYMENT.md](AWS_FARGATE_DEPLOYMENT.md) - Complete deployment guide
- [ARCHITECTURE.md](ARCHITECTURE.md) - Detailed architecture documentation
- [SETUP_GUIDE.md](SETUP_GUIDE.md) - Infrastructure setup guide

## 📊 Monitoring

### CloudWatch Logs

```bash
# View logs
aws logs tail /ecs/chat-app --follow

# Filter errors
aws logs filter-pattern /ecs/chat-app --filter-pattern "ERROR"
```

### Health Checks

```bash
# Check application health
curl http://your-alb-dns/health

# Check ECS service
aws ecs describe-services --cluster chat-app-cluster --services chat-app-service
```

### Metrics

- **Application**: Request count, error rate, response time
- **Infrastructure**: CPU/Memory utilization, network traffic
- **Database**: Connections, query performance

### Scaling

```bash
# Scale ECS service
aws ecs update-service \
  --cluster chat-app-cluster \
  --service chat-app-service \
  --desired-count 4
```

## 🔒 Security

### Best Practices Implemented

- ✅ **Secrets Management**: AWS Secrets Manager for credentials
- ✅ **Network Isolation**: Database in private subnets
- ✅ **Encryption**: Data encrypted at rest and in transit
- ✅ **JWT Authentication**: Secure token-based auth
- ✅ **Password Hashing**: bcrypt with cost factor 12
- ✅ **RBAC**: Role-based access control
- ✅ **CORS**: Configured for frontend domains
- ✅ **Security Groups**: Minimal port exposure
- ✅ **IAM**: Least privilege access

### Environment Variables

Never commit sensitive data. Use:
- `.env` for local development (gitignored)
- AWS Secrets Manager for production
- GitHub Secrets for CI/CD

## 🤝 Contributing

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit changes (`git commit -m 'Add amazing feature'`)
6. Push to branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Standards

- Follow Go best practices
- Write tests for new features
- Update documentation
- Use meaningful commit messages

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Abraham Nyagar** - [GitHub](https://github.com/Nyagar-Abraham)

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [Stream Chat](https://getstream.io/)
- [AWS](https://aws.amazon.com/)
- [Terraform](https://www.terraform.io/)

## 📞 Support

For issues and questions:
- Open an issue on GitHub
- Check existing documentation
- Review API documentation at `/swagger`

## 🗺️ Roadmap

- [ ] WebSocket support for real-time updates
- [ ] Message reactions and threading
- [ ] File upload and sharing
- [ ] Push notifications
- [ ] Advanced search functionality
- [ ] Analytics dashboard
- [ ] Rate limiting
- [ ] API versioning

---

**Built with ❤️ using Go and AWS**
