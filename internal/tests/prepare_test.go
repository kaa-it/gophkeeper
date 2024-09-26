package tests_test

import (
	"net"
	"testing"
	"time"

	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/creditcard"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/text"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kaa-it/gophkeeper/internal/client"
	"github.com/kaa-it/gophkeeper/internal/client/auth"
	"github.com/kaa-it/gophkeeper/internal/client/keeper"
	"github.com/kaa-it/gophkeeper/internal/pb"
	"github.com/kaa-it/gophkeeper/internal/server"
	serverAuth "github.com/kaa-it/gophkeeper/internal/server/application/grpc/auth"
	serverKeeper "github.com/kaa-it/gophkeeper/internal/server/application/grpc/keeper"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/credentials"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/file"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/user"
)

const (
	secretKey       = "secret"
	tokenDuration   = 15 * time.Minute
	refreshDuration = 30 * time.Second
)

func newTestClients(t *testing.T, address string) (*auth.Client, *keeper.Client) {
	cc1, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	authClient := auth.NewClient(cc1)
	interceptor := auth.NewAuthInterceptor(authClient, client.AuthMethods(), refreshDuration)

	cc2, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	require.NoError(t, err)

	keeperClient := keeper.NewClient(cc2)

	return authClient, keeperClient
}

func newTestSever(
	t *testing.T,
	userStore user.Repository,
	fileStore file.Repository,
	credentialsStore credentials.Repository,
	creditCardStore creditcard.Repository,
	textStore text.Repository,
) string {
	jwtManager := serverAuth.NewJWTManager(secretKey, tokenDuration)

	authServer := serverAuth.NewServer(userStore, jwtManager)
	keeperServer := serverKeeper.NewServer(fileStore, credentialsStore, creditCardStore, textStore)

	interceptor := serverAuth.NewAuthInterceptor(jwtManager, server.ProtectedMethods())

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterKeeperServiceServer(grpcServer, keeperServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go func() {
		_ = grpcServer.Serve(listener)
	}()

	return listener.Addr().String()
}
