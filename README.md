# Bookstore API

A RESTful and gRPC API for managing books, authors, and categories built with Go, Fiber, and PostgreSQL.

## Features

- **Books Management**: CRUD operations for books
- **Authors Management**: CRUD operations for authors  
- **Categories Management**: CRUD operations for categories
- **Dual API**: Both REST (Fiber) and gRPC endpoints
- **PostgreSQL**: Robust data persistence
- **Validation**: Input validation and error handling

## Project Structure

```
bookstore-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── models/
│   ├── handlers/
│   ├── services/
│   └── grpc/
├── proto/
├── migrations/
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

1. Install dependencies: `go mod tidy`
2. Set up PostgreSQL database
3. Run migrations
4. Start the server: `go run cmd/server/main.go`

## API Endpoints

### REST API (Fiber)
- `GET /api/v1/books` - List all books
- `POST /api/v1/books` - Create a new book
- `GET /api/v1/books/:id` - Get book by ID
- `PUT /api/v1/books/:id` - Update book
- `DELETE /api/v1/books/:id` - Delete book

Similar endpoints for authors and categories.

### gRPC
- BookService, AuthorService, CategoryService with full CRUD operations
