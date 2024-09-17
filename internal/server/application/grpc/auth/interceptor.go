package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	jwtManager      *JWTManager
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(jwtManager *JWTManager, accessibleRoles map[string][]string) *Interceptor {
	return &Interceptor{
		jwtManager:      jwtManager,
		accessibleRoles: accessibleRoles,
	}
}

func (interceptor *Interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *Interceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := interceptor.authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, ss)
	}
}

func (interceptor *Interceptor) authorize(ctx context.Context, _ string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata not found")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization not found")
	}

	accessToken := values[0]
	_, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return status.Errorf(codes.PermissionDenied, "permission denied")
}
