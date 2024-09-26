package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kaa-it/gophkeeper/internal/pb"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/text"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/user"
)

func TestClientText(t *testing.T) {
	t.Parallel()

	userStore := user.NewInMemoryUserStore()
	textStore := text.NewInMemoryTextStore()

	serverAddress := newTestSever(t, userStore, nil, nil, nil, textStore)

	authClient, keeperClient := newTestClients(t, serverAddress)

	err := authClient.Register("admin", "admin", "admin")
	require.NoError(t, err)

	err = authClient.Login("admin", "admin")
	require.NoError(t, err)

	textReq := &pb.Text{
		Metadata: "test@mail.ru",
		Text:     "test text",
	}

	_, err = keeperClient.UploadText(textReq)
	require.NoError(t, err)
}
