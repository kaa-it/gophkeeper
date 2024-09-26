package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kaa-it/gophkeeper/internal/pb"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/credentials"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/user"
)

func TestClientCredentials(t *testing.T) {
	t.Parallel()

	userStore := user.NewInMemoryUserStore()
	credentialsStore := credentials.NewInMemoryCredentialsStore()

	serverAddress := newTestSever(t, userStore, nil, credentialsStore, nil, nil)

	authClient, keeperClient := newTestClients(t, serverAddress)

	err := authClient.Register("admin", "admin", "admin")
	require.NoError(t, err)

	err = authClient.Login("admin", "admin")
	require.NoError(t, err)

	credentialsReq := &pb.Credentials{
		Metadata: "test@mail.ru",
		Login:    "test",
		Password: "8888",
	}

	_, err = keeperClient.UploadCredentials(credentialsReq)
	require.NoError(t, err)
}
