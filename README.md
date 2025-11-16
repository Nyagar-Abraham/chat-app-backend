# Chat Application API

A multi-tenant real-time chat application built with Go, featuring role-based access control (RBAC), JWT authentication, and Stream Chat integration for scalable messaging.

## 🚀 Features

- **Multi-Tenant Architecture** - Isolated organizations with tenant-based data segregation
- **Role-Based Access Control (RBAC)** - Four user roles: Admin, Moderator, Member, and Guest
- **JWT Authentication** - Secure token-based authentication
- **Real-Time Messaging** - Powered by Stream Chat API
- **RESTful API** - Clean and intuitive endpoints
- **PostgreSQL Database** - Reliable data persistence with GORM
- **Swagger Documentation** - Interactive API documentation
- **CORS Support** - Cross-origin resource sharing enabled

## 📋 Table of Contents

- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
- [User Roles & Permissions](#user-roles--permissions)

## 🛠 Technology Stack

- **Language**: Go 1.24
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Real-Time Chat**: Stream Chat Go SDK
- **Password Hashing**: bcrypt
- **API Documentation**: Swagger (gin-swagger)
- **Environment Management**: godotenv

## 📦 Prerequisites

- Go 1.24 or higher
- Docker & Docker Compose (recommended) OR PostgreSQL 12+
- Stream Chat account ([Get API credentials](https://getstream.io/chat/))

## 🔧 Installation

1. **Clone the repository**
```bash
git clone https://github.com/Nyagar-Abraham/chat-app.git
cd chat-app
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up PostgreSQL database**

**Option A: Using Docker Compose (Recommended)**
```bash
docker-compose up -d
```

**Option B: Local PostgreSQL**
```bash
createdb chat_db
```

## ⚙️ Configuration

1. **Copy the example environment file**
```bash
cp .env.example .env
```

2. **Configure environment variables in `.env`**
```env
DATABASE_URL=postgresql://chatuser:chatpassword@localhost:5432/chat_db?sslmode=disable
STREAM_API_KEY=<your_stream_api_key>
STREAM_API_SECRET=<your_stream_api_secret>
JWT_SECRET=<your_jwt_secret>
MIGRATE_DB=true
PORT=8085
```

**Environment Variables:**
- `DATABASE_URL` - PostgreSQL connection string
- `STREAM_API_KEY` - Stream Chat API key
- `STREAM_API_SECRET` - Stream Chat API secret
- `JWT_SECRET` - Secret key for JWT token generation
- `MIGRATE_DB` - Set to `true` to auto-migrate database schema (development only)
- `PORT` - Server port (default: 8085)

## 🚀 Running the Application

**Start the server**
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8085`

**Access Swagger Documentation**
```
http://localhost:8085/swagger/index.html
```

## 📁 Project Structure

```
chat-app/
├── cmd/
│   └── main.go              # Application entry point
├── db/
│   └── db.go                # Database connection and migration
├── handlers/
│   ├── auth.go              # Authentication handlers (login, register)
│   ├── channel.go           # Channel management handlers
│   ├── stream.go            # Stream Chat messaging handlers
│   └── tenant_user.go       # Tenant and user management handlers
├── middleware/
│   ├── auth.go              # Authentication middleware
│   ├── jwt.go               # JWT middleware
│   └── rbac.go              # Role-based access control middleware
├── models/
│   └── models.go            # Data models (User, Tenant, Channel)
├── services/
│   ├── bcrypt.go            # Password hashing service
│   ├── channel.go           # Channel service logic
│   └── stream.go            # Stream Chat integration
├── utils/
│   └── jwt.go               # JWT token utilities
├── .env.example             # Example environment configuration
├── go.mod                   # Go module dependencies
└── go.sum                   # Dependency checksums
```

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword",
  "role": "MEMBER",
  "org_name": "Acme Corp"
}
```

**Response:**
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "MEMBER",
  "token": "jwt_token",
  "tenant_id": "tenant_uuid"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "jwt_token",
  "message": "John Doe"
}
```

#### Get Current User
```http
GET /me
Authorization: Bearer <token>
```

### Tenant Endpoints

#### Create Tenant (Admin Only)
```http
POST /tenants
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "New Organization"
}
```

#### List All Tenants
```http
GET /tenants
```

#### Get Tenant by ID
```http
GET /tenants/:id
```

### User Management Endpoints

#### Create User (Admin/Moderator)
```http
POST /users
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "password": "password123",
  "role": "MEMBER",
  "tenant_id": "tenant_uuid"
}
```

#### List Users
```http
GET /users
Authorization: Bearer <token>
```

#### Update User (Admin/Moderator)
```http
PUT /users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Jane Doe",
  "role": "MODERATOR"
}
```

#### Delete User (Admin Only)
```http
DELETE /users/:id
Authorization: Bearer <token>
```

### Channel Endpoints

#### Create Channel (Admin/Moderator)
```http
POST /channels
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "General",
  "description": "General discussion channel"
}
```

**Note**: The channel creator is automatically added as a member.

#### List Channels
```http
GET /channels
Authorization: Bearer <token>
```

#### Add User to Channel (Admin/Moderator)
```http
POST /channels/:id/members
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": "user_uuid"
}
```

#### Remove User from Channel (Admin/Moderator)
```http
DELETE /channels/:id/members/:user_id
Authorization: Bearer <token>
```

#### Join Channel
```http
POST /channels/:id/join
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Joined channel successfully"
}
```

#### Leave Channel
```http
POST /channels/:id/leave
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Left channel successfully"
}
```

#### List Channel Members
```http
GET /channels/:id/members
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": "user_uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "MEMBER",
    "tenant_id": "tenant_uuid"
  }
]
```

### Messaging Endpoints

#### Get Stream Token
```http
GET /stream/token
Authorization: Bearer <token>
```

**Response:**
```json
{
  "token": "stream_chat_token"
}
```

#### Send Message
```http
POST /messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "stream_id": "channel_stream_id",
  "text": "Hello, World!"
}
```

**Note**: User must be a member of the channel to send messages. Returns 403 Forbidden if not a member.

#### Get Messages
```http
GET /messages/:stream_id
Authorization: Bearer <token>
```

**Note**: User must be a member of the channel to view messages. Returns 403 Forbidden if not a member.

## 👥 User Roles & Permissions

| Role | Permissions |
|------|-------------|
| **ADMIN** | Full access - manage tenants, users, channels, and messages |
| **MODERATOR** | Create/update users, create channels, send messages |
| **MEMBER** | Send messages, view channels |
| **GUEST** | Limited read access |

### Role-Based Endpoint Access

- **Tenant Creation**: Admin only
- **User Creation/Update**: Admin, Moderator
- **User Deletion**: Admin only
- **Channel Creation**: Admin, Moderator
- **Add/Remove Channel Members**: Admin, Moderator
- **Join/Leave Channels**: All authenticated users
- **Messaging**: Channel members only (tenant-isolated)

## 🔐 Authentication Flow

1. **Register** - Create account with organization name (auto-creates or joins tenant)
2. **Login** - Receive JWT token
3. **Authenticate** - Include token in `Authorization: Bearer <token>` header
4. **Access Resources** - Token contains user_id, tenant_id, and role claims

## 🔒 Channel Membership & Security

### Tenant Isolation
All channels are isolated by tenant. Users can only:
- View channels within their tenant
- Join channels within their tenant
- Send/receive messages in channels they are members of

### Membership Flow
1. **Admin/Moderator creates a channel** - Creator is automatically added as a member
2. **Users join channels** - Use `/channels/:id/join` endpoint
3. **Admin/Moderator adds users** - Use `/channels/:id/members` endpoint
4. **Send messages** - Only channel members can send/view messages

### Security Features
- ✅ Tenant-scoped channel access
- ✅ Membership validation for messaging
- ✅ Role-based member management
- ✅ Automatic Stream Chat synchronization

## 🧪 Testing

Use the Swagger UI at `/swagger/index.html` for interactive API testing, or use tools like:
- **Postman** - Import endpoints and test
- **cURL** - Command-line testing
- **HTTPie** - User-friendly HTTP client

**Example cURL request:**
```bash
curl -X POST http://localhost:8085/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

## 📝 License

This project is available for use under standard software licensing terms.

## 🤝 Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📧 Contact

For questions or support, please contact the development team.

---

**Built with ❤️ using Go and Stream Chat**
