# Database Schema Specification

## Overview
The book management system uses PostgreSQL with two main tables: `users` for authentication and `books` for catalog management.

## Tables

### users
User authentication and profile management table.

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
```

#### Fields Description
- `id`: Primary key, auto-increment
- `email`: Unique user email for login
- `password_hash`: Bcrypt hashed password
- `first_name`: User's first name
- `last_name`: User's last name
- `role`: User role (`admin` | `member`)
- `status`: Account status (`active` | `inactive`)
- `created_at`: Record creation timestamp (UTC)
- `updated_at`: Record last update timestamp (UTC)

### books
Book catalog and inventory management table.

```sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) UNIQUE,
    publisher VARCHAR(255),
    publication_year INTEGER,
    genre VARCHAR(100),
    description TEXT,
    pages INTEGER,
    language VARCHAR(50) DEFAULT 'English',
    price DECIMAL(10,2),
    quantity INTEGER NOT NULL DEFAULT 1,
    available_quantity INTEGER NOT NULL DEFAULT 1,
    location VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_books_title ON books(title);
CREATE INDEX idx_books_author ON books(author);
CREATE UNIQUE INDEX idx_books_isbn ON books(isbn) WHERE isbn IS NOT NULL;
CREATE INDEX idx_books_genre ON books(genre);
CREATE INDEX idx_books_status ON books(status);
```

#### Fields Description
- `id`: Primary key, auto-increment
- `title`: Book title (searchable)
- `author`: Author name (searchable)
- `isbn`: International Standard Book Number (unique)
- `publisher`: Publishing company
- `publication_year`: Year of publication
- `genre`: Book category/genre
- `description`: Book summary/description
- `pages`: Number of pages
- `language`: Book language (default: English)
- `price`: Book price (decimal)
- `quantity`: Total number of copies
- `available_quantity`: Currently available copies
- `location`: Physical location (shelf/section)
- `status`: Book availability status (`available` | `unavailable`)
- `created_at`: Record creation timestamp (UTC)
- `updated_at`: Record last update timestamp (UTC)

## Data Constraints

### Business Rules
1. **Email Uniqueness**: Each user must have a unique email address
2. **ISBN Uniqueness**: If provided, ISBN must be unique across all books
3. **Inventory Logic**: `available_quantity` should never exceed `quantity`
4. **Role Validation**: User role must be either 'admin' or 'member'
5. **Status Validation**: User status must be 'active' or 'inactive'
6. **Book Status**: Book status must be 'available' or 'unavailable'

### Required Fields
- **users**: email, password_hash, first_name, last_name
- **books**: title, author, quantity, available_quantity

### Default Values
- **users**: role='member', status='active'
- **books**: language='English', quantity=1, available_quantity=1, status='available'

## Environment Configuration

Database connection requires these environment variables:
```bash
BOOKMS_DB_HOST=localhost
BOOKMS_DB_PORT=5432
BOOKMS_DB_USER=bookms_user
BOOKMS_DB_PASSWORD=secure_password
BOOKMS_DB_NAME=book_management
BOOKMS_DB_MAX_OPEN_CONNS=25
BOOKMS_DB_MAX_IDLE_CONNS=5
BOOKMS_DB_CONN_MAX_LIFETIME=300
```

## Migration Notes
- All timestamps stored in UTC
- Password hashing uses bcrypt with cost 12
- Database connection pool configured via environment variables
- Indexes optimized for search operations on title, author, and email