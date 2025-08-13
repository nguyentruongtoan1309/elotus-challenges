# JWT Authentication & File Upload Server

A Go-based HTTP server that provides JWT authentication and secure file upload functionality. Built as a skills assessment project demonstrating backend development capabilities with Go.

## Features

### 1. JWT Authentication System

- **User Registration**: Create new user accounts with username/password
- **User Login**: Authenticate users and receive JWT tokens
- **Token Revocation**: Logout functionality that invalidates tokens
- **HS256 Signing**: Uses HMAC SHA-256 for token signing
- **Token Expiration**: 24-hour token lifetime with automatic cleanup

### 2. Secure File Upload API

- **Image Upload**: Accepts image files through multipart form data
- **Authorization Required**: All uploads require valid JWT tokens
- **File Validation**: Ensures uploaded files are images and under 8MB
- **Metadata Storage**: Stores file information and HTTP metadata in database
- **Temporary Storage**: Files saved to `/tmp` directory with unique names

## Quick Start

### Prerequisites

- Go 1.21 or higher
- SQLite3 (included with Go sqlite driver)

### Installation & Running

1. **Clone and setup the project:**

```bash
git clone <repository-url>
cd file-uploader
```

2. **Install dependencies:**

```bash
go mod tidy
```

3. **Run the server:**

```bash
go run main.go
```

The server will start on port 8080 by default. You can set a custom port using the `PORT` environment variable:

```bash
PORT=3000 go run main.go
```

4. **Access the test interface:**
   Open your browser and navigate to `http://localhost:8080` to access the simple HTML test interface.

### Running with Docker

1. **Build the Docker image:**

```bash
docker build -t file-uploader .
```

2. **Run the container:**

```bash
docker run -p 8080:8080 file-uploader
```

Or with custom environment variables:

```bash
docker run -p 8080:8080 -e JWT_SECRET="your-secure-secret" -e PORT=8080 file-uploader
```

3. **Using Docker Compose (recommended):**

```bash
docker-compose up
```

To run in detached mode:

```bash
docker-compose up -d
```

To stop the services:

```bash
docker-compose down
```

4. **Access the application:**
   Open your browser and navigate to `http://localhost:8080` to access the application.

## Environment Variables

| Variable     | Description        | Default           |
| ------------ | ------------------ | ----------------- |
| `PORT`       | Server port        | `8080`            |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |

**Important**: Change the JWT_SECRET in production:

```bash
export JWT_SECRET="your-super-secure-secret-key-here"
go run main.go
```

## API Documentation

### Authentication Endpoints

#### POST /api/v1/register

Register a new user account.

**Request Body:**

```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Response (201 Created):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "created_at": "2024-01-01T12:00:00Z"
  },
  "message": "User registered successfully"
}
```

#### POST /api/v1/login

Authenticate an existing user.

**Request Body:**

```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "created_at": "2024-01-01T12:00:00Z"
  },
  "message": "Login successful"
}
```

#### POST /api/v1/revoke

Revoke (logout) the current JWT token.

**Headers:**

```
Authorization: Bearer <your-jwt-token>
```

**Response (200 OK):**

```json
{
  "message": "Token revoked successfully"
}
```

### File Upload Endpoint

#### POST /api/v1/upload

Upload an image file with authentication.

**Headers:**

```
Authorization: Bearer <your-jwt-token>
```

**Form Data:**

- `data`: Image file (required)

**Response (201 Created):**

```json
{
  "message": "File uploaded successfully",
  "file_id": 1,
  "file_url": "/files/1",
  "public_url": "/public/files/1",
  "metadata": {
    "id": 1,
    "user_id": 1,
    "filename": "image.jpg",
    "content_type": "image/jpeg",
    "size": 1024000,
    "file_path": "/tmp/upload_1_1704110400_image.jpg",
    "user_agent": "Mozilla/5.0...",
    "remote_addr": "127.0.0.1:54321",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### Error Responses

All endpoints return JSON error responses with appropriate HTTP status codes:

```json
{
  "error": "Error description here"
}
```

Common status codes:

- `400 Bad Request`: Invalid input or file validation failed
- `401 Unauthorized`: Missing or invalid JWT token
- `409 Conflict`: Username already exists (registration)
- `500 Internal Server Error`: Server-side errors

## File Upload Validation

The upload endpoint enforces several security measures:

1. **Authentication Required**: Valid JWT token must be provided
2. **File Type Validation**: Only image files are accepted:
   - JPEG/JPG
   - PNG
   - GIF
   - WebP
   - BMP
   - TIFF
   - SVG
3. **File Size Limit**: Maximum 8MB per file
4. **Metadata Logging**: Captures and stores:
   - File information (name, size, type)
   - User information (from JWT)
   - HTTP metadata (User-Agent, IP address)
   - Upload timestamp

## Database Schema

The application uses SQLite with two main tables:

### Users Table

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,  -- bcrypt hashed
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Files Table

```sql
CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    filename TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    user_agent TEXT,
    remote_addr TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
```

## Project Structure

```
file-uploader/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── handlers/
│   ├── auth.go            # Authentication handlers
│   ├── static.go          # Serve static files handlers
│   └── upload.go          # File upload handlers
├── middleware/
│   └── auth.go            # JWT authentication
├── models/
│   ├── user.go            # User database model
│   └── file.go            # File metadata model
└── utils/
    ├── jwt.go             # JWT token utilities
    └── tokenblacklist.go  # Token revocation management
```

## Design Decisions

### Architecture Choices

1. **Modular Structure**: Separated concerns into handlers, models, middleware, and utilities
2. **SQLite Database**: Simple, file-based database perfect for development and testing
3. **In-Memory Token Blacklist**: Simple revocation mechanism with automatic cleanup
4. **Bcrypt Password Hashing**: Industry-standard password security
5. **Gorilla Mux Router**: Popular, feature-rich HTTP router for Go

### Security Considerations

1. **JWT with HS256**: Symmetric signing for simplicity while maintaining security
2. **Token Expiration**: 24-hour lifetime prevents indefinite access
3. **Password Validation**: Minimum length requirements and secure hashing
4. **File Validation**: Strict content-type and size checking
5. **IP Logging**: Tracks upload sources for security auditing

### Trade-offs Made

1. **In-Memory Blacklist**: Simple but not distributed-system friendly (would use Redis in production)
2. **SQLite**: Easy setup but not suitable for high-concurrency production use
3. **Temporary File Storage**: Simple but would use cloud storage in production
4. **Basic HTML Interface**: Functional but not production-ready UI

## Testing the Application

### Using cURL

1. **Register a user:**

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

2. **Login:**

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

3. **Upload a file:**

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -F "data=@/path/to/your/image.jpg"
```

4. **Revoke token:**

```bash
curl -X POST http://localhost:8080/api/v1/revoke \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### Using the Web Interface

Navigate to `http://localhost:8080` in your browser to use the simple HTML interface for testing all functionality.

## Production Considerations

For production deployment, consider:

1. **Environment Variables**: Set secure JWT_SECRET
2. **Database**: Migrate to PostgreSQL or MySQL
3. **File Storage**: Use cloud storage (AWS S3, Google Cloud Storage)
4. **Token Blacklist**: Use Redis for distributed token revocation
5. **HTTPS**: Enable TLS encryption
6. **Rate Limiting**: Implement request rate limiting
7. **Logging**: Add structured logging with log levels
8. **Monitoring**: Add health checks and metrics
9. **Docker**: Containerize the application
10. **Load Balancing**: Use reverse proxy (nginx) for multiple instances
