package auth

import (
	"context"

	"github.com/kaa-it/gophkeeper/internal/server/domain"

	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/user"

	"github.com/kaa-it/gophkeeper/internal/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
func (s *Server) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if (req.GetLogin() == "") || (req.GetPassword() == "") {
		return nil, status.Errorf(codes.InvalidArgument, "invalid login or password")
	}

	foundedUser, err := s.userStore.Get(req.GetLogin())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to find foundedUser: %v", err)
	}

	if !foundedUser.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid login or password")
	}

	token, err := s.jwtManager.Generate(foundedUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	res := &pb.LoginResponse{
		AccessToken: token,
	}

	return res, nil
}

// Register handles the user registration process by validating input data and storing a new user in the repository.
func (s *Server) Register(_ context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	if (req.GetUsername() == "") || (req.GetLogin() == "") || (req.GetPassword() == "") {
		return nil, status.Errorf(codes.InvalidArgument, "invalid username, login or password")
	}

	newUser, err := domain.NewUser(req.GetUsername(), req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "password too long: %v", err)
	}

	if err := s.userStore.Save(newUser); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "failed to save user: %v", err)
	}

	return &emptypb.Empty{}, nil
}
