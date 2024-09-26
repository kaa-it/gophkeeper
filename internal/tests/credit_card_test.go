package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kaa-it/gophkeeper/internal/pb"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/creditcard"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/user"
)

func TestClientCreditCard(t *testing.T) {
	t.Parallel()

	userStore := user.NewInMemoryUserStore()
	creditCardStore := creditcard.NewInMemoryCreditCardStore()

	serverAddress := newTestSever(t, userStore, nil, nil, creditCardStore, nil)

	authClient, keeperClient := newTestClients(t, serverAddress)

	err := authClient.Register("admin", "admin", "admin")
	require.NoError(t, err)

	err = authClient.Login("admin", "admin")
	require.NoError(t, err)

	creditCardReq := &pb.CreditCard{
		Metadata: "T-Bank",
		Name:     "Andrey Kuznetsov",
		Month:    "02",
		Year:     "25",
		Number:   "1234 5678 9012 3456",
		BillingAddress: &pb.BillingAddress{
			Country:  "Russia",
			Street:   "Tverskaya",
			City:     "Moscow",
			Postcode: "111456",
		},
	}

	_, err = keeperClient.UploadCreditCard(creditCardReq)
	require.NoError(t, err)
}
