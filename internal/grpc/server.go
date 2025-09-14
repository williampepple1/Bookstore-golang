package grpc

import (
	"bookstore-api/internal/config"
	"bookstore-api/internal/services"
	pb "bookstore-api/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

// GRPCServer represents the gRPC server
type GRPCServer struct {
	pb.UnimplementedAuthorServiceServer
	pb.UnimplementedCategoryServiceServer
	pb.UnimplementedBookServiceServer
	pb.UnimplementedHealthServiceServer

	authorService   *services.AuthorService
	categoryService *services.CategoryService
	bookService     *services.BookService
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{
		authorService:   services.NewAuthorService(),
		categoryService: services.NewCategoryService(),
		bookService:     services.NewBookService(),
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start(cfg *config.Config) error {
	lis, err := net.Listen("tcp", cfg.GRPC.Host+":"+cfg.GRPC.Port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	// Register services
	pb.RegisterAuthorServiceServer(grpcServer, s)
	pb.RegisterCategoryServiceServer(grpcServer, s)
	pb.RegisterBookServiceServer(grpcServer, s)
	pb.RegisterHealthServiceServer(grpcServer, s)

	log.Printf("Starting gRPC server on %s:%s", cfg.GRPC.Host, cfg.GRPC.Port)
	return grpcServer.Serve(lis)
}

// Health Check implementation
func (s *GRPCServer) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status:  pb.HealthCheckResponse_SERVING,
		Message: "gRPC service is healthy",
	}, nil
}
