package auth

import (
	"context"

	"github.com/kaa-it/gophkeeper/internal/server/domain"

	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/user"

	"github.com/kaa-it/gophkeeper/internal/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is responsible for handling authentication-related requests.
type Server struct {
	userStore  user.Repository
	jwtManager *JWTManager
	pb.UnimplementedAuthServiceServer
}

// NewServer initializes and returns a new AuthServer with provided UserStore and JWTManager.
func NewServer(userStore user.Repository, jwtManager *JWTManager) *Server {
	return &Server{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

// Login handles the authentication process by validating user credentials and generating a JWT token if successful.
func (s *Server) Login(_ context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, err := s.userStore.GetUser(req.GetLogin())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	res := &pb.AuthResponse{
		AccessToken: token,
	}

	return res, nil
}

func (s *Server) Register(_ context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	user, err := domain.NewUser(req.GetUsername(), req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	if err := s.userStore.Save(user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save user: %v", err)
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	res := &pb.AuthResponse{
		AccessToken: token,
	}

	return res, nil
}
