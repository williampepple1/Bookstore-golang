package grpc

import (
	"bookstore-api/internal/models"
	pb "bookstore-api/proto"
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateBook implements the CreateBook gRPC method
func (s *GRPCServer) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	authorID, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return &pb.CreateBookResponse{
			Success: false,
			Message: "Invalid author ID",
		}, status.Error(codes.InvalidArgument, "Invalid author ID")
	}

	categoryID, err := uuid.Parse(req.CategoryId)
	if err != nil {
		return &pb.CreateBookResponse{
			Success: false,
			Message: "Invalid category ID",
		}, status.Error(codes.InvalidArgument, "Invalid category ID")
	}

	var publishedAt *time.Time
	if req.PublishedAt != "" {
		if parsed, err := time.Parse("2006-01-02T15:04:05Z07:00", req.PublishedAt); err == nil {
			publishedAt = &parsed
		}
	}

	book := &models.Book{
		Title:       req.Title,
		ISBN:        req.Isbn,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		PublishedAt: publishedAt,
		AuthorID:    authorID,
		CategoryID:  categoryID,
	}

	if err := s.bookService.CreateBook(book); err != nil {
		return &pb.CreateBookResponse{
			Success: false,
			Message: "Failed to create book: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateBookResponse{
		Success: true,
		Message: "Book created successfully",
		Book:    convertBookToProto(book),
	}, nil
}

// GetBook implements the GetBook gRPC method
func (s *GRPCServer) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.GetBookResponse{
			Success: false,
			Message: "Invalid book ID",
		}, status.Error(codes.InvalidArgument, "Invalid book ID")
	}

	book, err := s.bookService.GetBookByID(id)
	if err != nil {
		if err.Error() == "book not found" {
			return &pb.GetBookResponse{
				Success: false,
				Message: "Book not found",
			}, status.Error(codes.NotFound, "Book not found")
		}
		return &pb.GetBookResponse{
			Success: false,
			Message: "Failed to get book: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetBookResponse{
		Success: true,
		Message: "Book retrieved successfully",
		Book:    convertBookToProto(book),
	}, nil
}

// GetAllBooks implements the GetAllBooks gRPC method
func (s *GRPCServer) GetAllBooks(ctx context.Context, req *pb.GetAllBooksRequest) (*pb.GetAllBooksResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	books, total, err := s.bookService.GetAllBooks(page, limit)
	if err != nil {
		return &pb.GetAllBooksResponse{
			Success: false,
			Message: "Failed to get books: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	var protoBooks []*pb.Book
	for _, book := range books {
		protoBooks = append(protoBooks, convertBookToProto(&book))
	}

	return &pb.GetAllBooksResponse{
		Success: true,
		Message: "Books retrieved successfully",
		Books:   protoBooks,
		Pagination: &pb.Pagination{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      total,
			TotalPages: (total + int64(limit) - 1) / int64(limit),
		},
	}, nil
}

// UpdateBook implements the UpdateBook gRPC method
func (s *GRPCServer) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.UpdateBookResponse{
			Success: false,
			Message: "Invalid book ID",
		}, status.Error(codes.InvalidArgument, "Invalid book ID")
	}

	updates := &models.Book{
		Title:       req.Title,
		ISBN:        req.Isbn,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
	}

	// Parse optional fields
	if req.AuthorId != "" {
		authorID, err := uuid.Parse(req.AuthorId)
		if err != nil {
			return &pb.UpdateBookResponse{
				Success: false,
				Message: "Invalid author ID",
			}, status.Error(codes.InvalidArgument, "Invalid author ID")
		}
		updates.AuthorID = authorID
	}

	if req.CategoryId != "" {
		categoryID, err := uuid.Parse(req.CategoryId)
		if err != nil {
			return &pb.UpdateBookResponse{
				Success: false,
				Message: "Invalid category ID",
			}, status.Error(codes.InvalidArgument, "Invalid category ID")
		}
		updates.CategoryID = categoryID
	}

	if req.PublishedAt != "" {
		if parsed, err := time.Parse("2006-01-02T15:04:05Z07:00", req.PublishedAt); err == nil {
			updates.PublishedAt = &parsed
		}
	}

	if err := s.bookService.UpdateBook(id, updates); err != nil {
		if err.Error() == "book not found" {
			return &pb.UpdateBookResponse{
				Success: false,
				Message: "Book not found",
			}, status.Error(codes.NotFound, "Book not found")
		}
		return &pb.UpdateBookResponse{
			Success: false,
			Message: "Failed to update book: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateBookResponse{
		Success: true,
		Message: "Book updated successfully",
	}, nil
}

// DeleteBook implements the DeleteBook gRPC method
func (s *GRPCServer) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.DeleteBookResponse{
			Success: false,
			Message: "Invalid book ID",
		}, status.Error(codes.InvalidArgument, "Invalid book ID")
	}

	if err := s.bookService.DeleteBook(id); err != nil {
		if err.Error() == "book not found" {
			return &pb.DeleteBookResponse{
				Success: false,
				Message: "Book not found",
			}, status.Error(codes.NotFound, "Book not found")
		}
		return &pb.DeleteBookResponse{
			Success: false,
			Message: "Failed to delete book: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteBookResponse{
		Success: true,
		Message: "Book deleted successfully",
	}, nil
}

// SearchBooks implements the SearchBooks gRPC method
func (s *GRPCServer) SearchBooks(ctx context.Context, req *pb.SearchBooksRequest) (*pb.SearchBooksResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	books, total, err := s.bookService.SearchBooks(req.Query, page, limit)
	if err != nil {
		return &pb.SearchBooksResponse{
			Success: false,
			Message: "Failed to search books: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	var protoBooks []*pb.Book
	for _, book := range books {
		protoBooks = append(protoBooks, convertBookToProto(&book))
	}

	return &pb.SearchBooksResponse{
		Success: true,
		Message: "Books found successfully",
		Books:   protoBooks,
		Pagination: &pb.Pagination{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      total,
			TotalPages: (total + int64(limit) - 1) / int64(limit),
		},
	}, nil
}

// GetBooksByAuthor implements the GetBooksByAuthor gRPC method
func (s *GRPCServer) GetBooksByAuthor(ctx context.Context, req *pb.GetBooksByAuthorRequest) (*pb.GetBooksByAuthorResponse, error) {
	authorID, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return &pb.GetBooksByAuthorResponse{
			Success: false,
			Message: "Invalid author ID",
		}, status.Error(codes.InvalidArgument, "Invalid author ID")
	}

	page := int(req.Page)
	limit := int(req.Limit)
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	books, total, err := s.bookService.GetBooksByAuthor(authorID, page, limit)
	if err != nil {
		return &pb.GetBooksByAuthorResponse{
			Success: false,
			Message: "Failed to get books by author: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	var protoBooks []*pb.Book
	for _, book := range books {
		protoBooks = append(protoBooks, convertBookToProto(&book))
	}

	return &pb.GetBooksByAuthorResponse{
		Success: true,
		Message: "Books retrieved successfully",
		Books:   protoBooks,
		Pagination: &pb.Pagination{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      total,
			TotalPages: (total + int64(limit) - 1) / int64(limit),
		},
	}, nil
}

// GetBooksByCategory implements the GetBooksByCategory gRPC method
func (s *GRPCServer) GetBooksByCategory(ctx context.Context, req *pb.GetBooksByCategoryRequest) (*pb.GetBooksByCategoryResponse, error) {
	categoryID, err := uuid.Parse(req.CategoryId)
	if err != nil {
		return &pb.GetBooksByCategoryResponse{
			Success: false,
			Message: "Invalid category ID",
		}, status.Error(codes.InvalidArgument, "Invalid category ID")
	}

	page := int(req.Page)
	limit := int(req.Limit)
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	books, total, err := s.bookService.GetBooksByCategory(categoryID, page, limit)
	if err != nil {
		return &pb.GetBooksByCategoryResponse{
			Success: false,
			Message: "Failed to get books by category: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	var protoBooks []*pb.Book
	for _, book := range books {
		protoBooks = append(protoBooks, convertBookToProto(&book))
	}

	return &pb.GetBooksByCategoryResponse{
		Success: true,
		Message: "Books retrieved successfully",
		Books:   protoBooks,
		Pagination: &pb.Pagination{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      total,
			TotalPages: (total + int64(limit) - 1) / int64(limit),
		},
	}, nil
}

// UpdateBookStock implements the UpdateBookStock gRPC method
func (s *GRPCServer) UpdateBookStock(ctx context.Context, req *pb.UpdateBookStockRequest) (*pb.UpdateBookStockResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.UpdateBookStockResponse{
			Success: false,
			Message: "Invalid book ID",
		}, status.Error(codes.InvalidArgument, "Invalid book ID")
	}

	if err := s.bookService.UpdateBookStock(id, int(req.Stock)); err != nil {
		if err.Error() == "book not found" {
			return &pb.UpdateBookStockResponse{
				Success: false,
				Message: "Book not found",
			}, status.Error(codes.NotFound, "Book not found")
		}
		return &pb.UpdateBookStockResponse{
			Success: false,
			Message: "Failed to update book stock: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateBookStockResponse{
		Success: true,
		Message: "Book stock updated successfully",
	}, nil
}

// convertBookToProto converts a models.Book to pb.Book
func convertBookToProto(book *models.Book) *pb.Book {
	protoBook := &pb.Book{
		Id:          book.ID.String(),
		Title:       book.Title,
		Isbn:        book.ISBN,
		Description: book.Description,
		Price:       book.Price,
		Stock:       int32(book.Stock),
		CreatedAt:   book.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   book.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		AuthorId:    book.AuthorID.String(),
		CategoryId:  book.CategoryID.String(),
	}

	if book.PublishedAt != nil {
		protoBook.PublishedAt = book.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	// Convert author if it exists
	if book.Author.ID != uuid.Nil {
		protoBook.Author = convertAuthorToProto(&book.Author)
	}

	// Convert category if it exists
	if book.Category.ID != uuid.Nil {
		protoBook.Category = convertCategoryToProto(&book.Category)
	}

	return protoBook
}
