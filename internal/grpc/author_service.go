package grpc

import (
	"bookstore-api/internal/models"
	pb "bookstore-api/proto"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateAuthor implements the CreateAuthor gRPC method
func (s *GRPCServer) CreateAuthor(ctx context.Context, req *pb.CreateAuthorRequest) (*pb.CreateAuthorResponse, error) {
	author := &models.Author{
		Name:      req.Name,
		Email:     req.Email,
		Biography: req.Biography,
	}

	if err := s.authorService.CreateAuthor(author); err != nil {
		return &pb.CreateAuthorResponse{
			Success: false,
			Message: "Failed to create author: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateAuthorResponse{
		Success: true,
		Message: "Author created successfully",
		Author:  convertAuthorToProto(author),
	}, nil
}

// GetAuthor implements the GetAuthor gRPC method
func (s *GRPCServer) GetAuthor(ctx context.Context, req *pb.GetAuthorRequest) (*pb.GetAuthorResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.GetAuthorResponse{
			Success: false,
			Message: "Invalid author ID",
		}, status.Error(codes.InvalidArgument, "Invalid author ID")
	}

	author, err := s.authorService.GetAuthorByID(id)
	if err != nil {
		if err.Error() == "author not found" {
			return &pb.GetAuthorResponse{
				Success: false,
				Message: "Author not found",
			}, status.Error(codes.NotFound, "Author not found")
		}
		return &pb.GetAuthorResponse{
			Success: false,
			Message: "Failed to get author: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetAuthorResponse{
		Success: true,
		Message: "Author retrieved successfully",
		Author:  convertAuthorToProto(author),
	}, nil
}

// GetAllAuthors implements the GetAllAuthors gRPC method
func (s *GRPCServer) GetAllAuthors(ctx context.Context, req *pb.GetAllAuthorsRequest) (*pb.GetAllAuthorsResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	authors, total, err := s.authorService.GetAllAuthors(page, limit)
	if err != nil {
		return &pb.GetAllAuthorsResponse{
			Success: false,
			Message: "Failed to get authors: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	var protoAuthors []*pb.Author
	for _, author := range authors {
		protoAuthors = append(protoAuthors, convertAuthorToProto(&author))
	}

	return &pb.GetAllAuthorsResponse{
		Success: true,
		Message: "Authors retrieved successfully",
		Authors: protoAuthors,
		Pagination: &pb.Pagination{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      total,
			TotalPages: (total + int64(limit) - 1) / int64(limit),
		},
	}, nil
}

// UpdateAuthor implements the UpdateAuthor gRPC method
func (s *GRPCServer) UpdateAuthor(ctx context.Context, req *pb.UpdateAuthorRequest) (*pb.UpdateAuthorResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.UpdateAuthorResponse{
			Success: false,
			Message: "Invalid author ID",
		}, status.Error(codes.InvalidArgument, "Invalid author ID")
	}

	updates := &models.Author{
		Name:      req.Name,
		Email:     req.Email,
		Biography: req.Biography,
	}

	if err := s.authorService.UpdateAuthor(id, updates); err != nil {
		if err.Error() == "author not found" {
			return &pb.UpdateAuthorResponse{
				Success: false,
				Message: "Author not found",
			}, status.Error(codes.NotFound, "Author not found")
		}
		return &pb.UpdateAuthorResponse{
			Success: false,
			Message: "Failed to update author: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateAuthorResponse{
		Success: true,
		Message: "Author updated successfully",
	}, nil
}

// DeleteAuthor implements the DeleteAuthor gRPC method
func (s *GRPCServer) DeleteAuthor(ctx context.Context, req *pb.DeleteAuthorRequest) (*pb.DeleteAuthorResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.DeleteAuthorResponse{
			Success: false,
			Message: "Invalid author ID",
		}, status.Error(codes.InvalidArgument, "Invalid author ID")
	}

	if err := s.authorService.DeleteAuthor(id); err != nil {
		if err.Error() == "author not found" {
			return &pb.DeleteAuthorResponse{
				Success: false,
				Message: "Author not found",
			}, status.Error(codes.NotFound, "Author not found")
		}
		return &pb.DeleteAuthorResponse{
			Success: false,
			Message: "Failed to delete author: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteAuthorResponse{
		Success: true,
		Message: "Author deleted successfully",
	}, nil
}

// SearchAuthors implements the SearchAuthors gRPC method
func (s *GRPCServer) SearchAuthors(ctx context.Context, req *pb.SearchAuthorsRequest) (*pb.SearchAuthorsResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	authors, total, err := s.authorService.SearchAuthors(req.Query, page, limit)
	if err != nil {
		return &pb.SearchAuthorsResponse{
			Success: false,
			Message: "Failed to search authors: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	var protoAuthors []*pb.Author
	for _, author := range authors {
		protoAuthors = append(protoAuthors, convertAuthorToProto(&author))
	}

	return &pb.SearchAuthorsResponse{
		Success: true,
		Message: "Authors found successfully",
		Authors: protoAuthors,
		Pagination: &pb.Pagination{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      total,
			TotalPages: (total + int64(limit) - 1) / int64(limit),
		},
	}, nil
}

// convertAuthorToProto converts a models.Author to pb.Author
func convertAuthorToProto(author *models.Author) *pb.Author {
	protoAuthor := &pb.Author{
		Id:        author.ID.String(),
		Name:      author.Name,
		Email:     author.Email,
		Biography: author.Biography,
		CreatedAt: author.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: author.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Convert books if they exist
	for _, book := range author.Books {
		protoAuthor.Books = append(protoAuthor.Books, convertBookToProto(&book))
	}

	return protoAuthor
}
