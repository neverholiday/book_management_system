# API Specification

## Overview
RESTful API for book management system with JWT-based authentication. All endpoints return JSON responses and follow standard HTTP status codes.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

## Response Format

### Success Response (200 OK)
```json
{
  "data": <response_data>,
  "message": "<message from API>"
}
```

### Error Response (4xx, 5xx)
```json
{
  "message": "<response or some error message>"
}
```

## System Endpoints

### Health Check
```http
GET /healthz
```

**Response (200):**
```json
{
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z",
    "version": "1.0.0"
  },
  "message": "Service is healthy"
}
```

## Authentication Endpoints

### Register User
```http
POST /auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (200):**
```json
{
  "data": {
    "user": {
      "id": "user_12345",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "member",
      "status": "active"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "message": "User registered successfully"
}
```

### Login User
```http
POST /auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response (200):**
```json
{
  "data": {
    "user": {
      "id": "user_12345",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "member",
      "status": "active"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-15T12:00:00Z"
  },
  "message": "Login successful"
}
```

### Refresh Token
```http
POST /auth/refresh
```
**Headers:** `Authorization: Bearer <current_token>`

**Response (200):**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-16T12:00:00Z"
  },
  "message": "Token refreshed successfully"
}
```

### Get Profile
```http
GET /auth/profile
```
**Headers:** `Authorization: Bearer <jwt_token>`

**Response (200):**
```json
{
  "data": {
    "id": "user_12345",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "member",
    "status": "active",
    "created_date": "2024-01-01T12:00:00Z",
    "updated_date": "2024-01-01T12:00:00Z"
  },
  "message": "Profile retrieved successfully"
}
```

### Update Profile
```http
PUT /auth/profile
```
**Headers:** `Authorization: Bearer <jwt_token>`

**Request Body:**
```json
{
  "first_name": "Jane",
  "last_name": "Smith"
}
```

### Logout
```http
POST /auth/logout
```
**Headers:** `Authorization: Bearer <jwt_token>`

**Response (200):**
```json
{
  "data": null,
  "message": "Logged out successfully"
}
```

## User Management Endpoints
**Admin Only - Requires JWT token with admin role**

### Create User
```http
POST /users
```
**Headers:** `Authorization: Bearer <admin_jwt_token>`

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "password": "securepassword",
  "first_name": "New",
  "last_name": "User",
  "role": "member"
}
```

### Get All Users
```http
GET /users?limit=20&offset=0&role=member&status=active
```
**Headers:** `Authorization: Bearer <admin_jwt_token>`

**Query Parameters:**
- `limit` (optional): Number of records to return (default: 20)
- `offset` (optional): Number of records to skip (default: 0)
- `role` (optional): Filter by role (admin/member)
- `status` (optional): Filter by status (active/inactive)

**Response (200):**
```json
{
  "data": {
    "users": [
      {
        "id": "user_12345",
        "email": "user@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "role": "member",
        "status": "active",
        "created_date": "2024-01-01T12:00:00Z",
        "updated_date": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 1,
    "limit": 20,
    "offset": 0
  },
  "message": "Users retrieved successfully"
}
```

### Get User by ID
```http
GET /users/:id
```
**Headers:** `Authorization: Bearer <admin_jwt_token>`

### Update User
```http
PUT /users/:id
```
**Headers:** `Authorization: Bearer <admin_jwt_token>`

**Request Body:**
```json
{
  "first_name": "Updated",
  "last_name": "Name",
  "role": "admin",
  "status": "active"
}
```

### Delete User
```http
DELETE /users/:id
```
**Headers:** `Authorization: Bearer <admin_jwt_token>`

## Book Management Endpoints

### Get All Books (Public)
```http
GET /books?limit=20&offset=0&title=&author=&genre=&isbn=
```

**Query Parameters:**
- `limit` (optional): Number of records to return (default: 20)
- `offset` (optional): Number of records to skip (default: 0)
- `title` (optional): Search by title (partial match)
- `author` (optional): Search by author (partial match)
- `genre` (optional): Filter by genre
- `isbn` (optional): Search by ISBN

**Response (200):**
```json
{
  "data": {
    "books": [
      {
        "id": "book_67890",
        "title": "The Go Programming Language",
        "author": "Alan Donovan, Brian Kernighan",
        "isbn": "978-0134190440",
        "publisher": "Addison-Wesley",
        "publication_year": 2015,
        "genre": "Programming",
        "description": "The authoritative resource to writing clear and idiomatic Go",
        "pages": 380,
        "language": "English",
        "price": 45.99,
        "quantity": 5,
        "available_quantity": 3,
        "location": "Shelf A-1",
        "status": "available",
        "created_date": "2024-01-01T12:00:00Z",
        "updated_date": "2024-01-01T12:00:00Z"
      }
    ],
    "total": 1,
    "limit": 20,
    "offset": 0
  },
  "message": "Books retrieved successfully"
}
```

### Get Book by ID (Public)
```http
GET /books/:id
```

### Create Book (Admin Only)
```http
POST /books
```
**Headers:** `Authorization: Bearer <admin_jwt_token>`

**Request Body:**
```json
{
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "isbn": "978-0132350884",
  "publisher": "Prentice Hall",
  "publication_year": 2008,
  "genre": "Programming",
  "description": "A Handbook of Agile Software Craftsmanship",
  "pages": 464,
  "language": "English",
  "price": 42.99,
  "quantity": 3,
  "location": "Shelf B-2"
}
```

### Update Book (Admin Only)
```http
PUT /books/:id
```
**Headers:** `Authorization: Bearer <admin_jwt_token>`

**Request Body:** (any fields to update)
```json
{
  "available_quantity": 2,
  "location": "Shelf B-3"
}
```

### Delete Book (Admin Only)
```http
DELETE /books/:id
```
**Headers:** `Authorization: Bearer <admin_jwt_token>`

## HTTP Status Codes

- `200 OK`: Successful GET, PUT operations
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Insufficient permissions (not admin)
- `404 Not Found`: Resource not found
- `409 Conflict`: Duplicate resource (email, ISBN)
- `422 Unprocessable Entity`: Validation errors
- `500 Internal Server Error`: Server error

## Error Codes

- `INVALID_CREDENTIALS`: Login failed
- `EMAIL_ALREADY_EXISTS`: Email already registered
- `ISBN_ALREADY_EXISTS`: ISBN already exists
- `USER_NOT_FOUND`: User not found
- `BOOK_NOT_FOUND`: Book not found
- `INVALID_TOKEN`: JWT token is invalid
- `TOKEN_EXPIRED`: JWT token has expired
- `INSUFFICIENT_PERMISSIONS`: User lacks required permissions
- `VALIDATION_ERROR`: Request validation failed

## Rate Limiting
- **Authentication endpoints**: 5 requests per minute per IP
- **User management**: 100 requests per minute per user
- **Book endpoints**: 200 requests per minute per user

## JWT Token Configuration
- **Algorithm**: HS256
- **Expiry**: 24 hours (configurable via `BOOKMS_JWT_EXPIRY_HOURS`)
- **Refresh**: 7 days (configurable via `BOOKMS_JWT_REFRESH_EXPIRY_HOURS`)
- **Claims**: user_id, email, role, iat, exp
