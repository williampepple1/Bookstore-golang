package grpc

import (
	"bookstore-api/internal/models"
	pb "bookstore-api/proto"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateCategory implements the CreateCategory gRPC method
func (s *GRPCServer) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryService.CreateCategory(category); err != nil {
		return &pb.CreateCategoryResponse{
			Success: false,
			Message: "Failed to create category: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateCategoryResponse{
		Success:  true,
		Message:  "Category created successfully",
		Category: convertCategoryToProto(category),
	}, nil
}

// GetCategory implements the GetCategory gRPC method
func (s *GRPCServer) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.GetCategoryResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.GetCategoryResponse{
			Success: false,
			Message: "Invalid category ID",
		}, status.Error(codes.InvalidArgument, "Invalid category ID")
	}

	category, err := s.categoryService.GetCategoryByID(id)
	if err != nil {
		if err.Error() == "category not found" {
			return &pb.GetCategoryResponse{
				Success: false,
				Message: "Category not found",
			}, status.Error(codes.NotFound, "Category not found")
		}
		return &pb.GetCategoryResponse{
			Success: false,
			Message: "Failed to get category: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetCategoryResponse{
		Success:  true,
		Message:  "Category retrieved successfully",
		Category: convertCategoryToProto(category),
	}, nil
}

// GetAllCategories implements the GetAllCategories gRPC method
func (s *GRPCServer) GetAllCategories(ctx context.Context, req *pb.GetAllCategoriesRequest) (*pb.GetAllCategoriesResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	categories, total, err := s.categoryService.GetAllCategories(page, limit)
	if err != nil {
		return &pb.GetAllCategoriesResponse{
			Success: false,
			Message: "Failed to get categories: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	var protoCategories []*pb.Category
	for _, category := range categories {
		protoCategories = append(protoCategories, convertCategoryToProto(&category))
	}

	return &pb.GetAllCategoriesResponse{
		Success:    true,
		Message:    "Categories retrieved successfully",
		Categories: protoCategories,
		Pagination: &pb.Pagination{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      total,
			TotalPages: (total + int64(limit) - 1) / int64(limit),
		},
	}, nil
}

// UpdateCategory implements the UpdateCategory gRPC method
func (s *GRPCServer) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.UpdateCategoryResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.UpdateCategoryResponse{
			Success: false,
			Message: "Invalid category ID",
		}, status.Error(codes.InvalidArgument, "Invalid category ID")
	}

	updates := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryService.UpdateCategory(id, updates); err != nil {
		if err.Error() == "category not found" {
			return &pb.UpdateCategoryResponse{
				Success: false,
				Message: "Category not found",
			}, status.Error(codes.NotFound, "Category not found")
		}
		return &pb.UpdateCategoryResponse{
			Success: false,
			Message: "Failed to update category: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateCategoryResponse{
		Success: true,
		Message: "Category updated successfully",
	}, nil
}

// DeleteCategory implements the DeleteCategory gRPC method
func (s *GRPCServer) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.DeleteCategoryResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.DeleteCategoryResponse{
			Success: false,
			Message: "Invalid category ID",
		}, status.Error(codes.InvalidArgument, "Invalid category ID")
	}

	if err := s.categoryService.DeleteCategory(id); err != nil {
		if err.Error() == "category not found" {
			return &pb.DeleteCategoryResponse{
				Success: false,
				Message: "Category not found",
			}, status.Error(codes.NotFound, "Category not found")
		}
		return &pb.DeleteCategoryResponse{
			Success: false,
			Message: "Failed to delete category: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteCategoryResponse{
		Success: true,
		Message: "Category deleted successfully",
	}, nil
}

// SearchCategories implements the SearchCategories gRPC method
func (s *GRPCServer) SearchCategories(ctx context.Context, req *pb.SearchCategoriesRequest) (*pb.SearchCategoriesResponse, error) {
	page := int(req.Page)
	limit := int(req.Limit)
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	categories, total, err := s.categoryService.SearchCategories(req.Query, page, limit)
	if err != nil {
		return &pb.SearchCategoriesResponse{
			Success: false,
			Message: "Failed to search categories: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	var protoCategories []*pb.Category
	for _, category := range categories {
		protoCategories = append(protoCategories, convertCategoryToProto(&category))
	}

	return &pb.SearchCategoriesResponse{
		Success:    true,
		Message:    "Categories found successfully",
		Categories: protoCategories,
		Pagination: &pb.Pagination{
			Page:       int32(page),
			Limit:      int32(limit),
			Total:      total,
			TotalPages: (total + int64(limit) - 1) / int64(limit),
		},
	}, nil
}

// convertCategoryToProto converts a models.Category to pb.Category
func convertCategoryToProto(category *models.Category) *pb.Category {
	protoCategory := &pb.Category{
		Id:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Convert books if they exist
	for _, book := range category.Books {
		protoCategory.Books = append(protoCategory.Books, convertBookToProto(&book))
	}

	return protoCategory
}
