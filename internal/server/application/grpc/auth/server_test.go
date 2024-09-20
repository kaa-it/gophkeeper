package auth_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kaa-it/gophkeeper/internal/server/application/grpc/auth"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kaa-it/gophkeeper/internal/pb"
	"github.com/kaa-it/gophkeeper/internal/server/domain"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/user"
)

func TestServerRegister(t *testing.T) {
	t.Parallel()

	jwtManager := auth.NewJWTManager("secretKey", 15*time.Minute)

	emptyStore := user.NewInMemoryUserStore()

	storeWithUser := user.NewInMemoryUserStore()
	newUser, err := domain.NewUser("Andrey", "admin", "test")
	require.NoError(t, err)

	err = storeWithUser.Save(newUser)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		username string
		login    string
		password string
		store    user.Repository
		code     codes.Code
	}{
		{
			name:     "should return error if username is empty",
			username: "",
			login:    "admin",
			password: "test",
			store:    emptyStore,
			code:     codes.InvalidArgument,
		},
		{
			name:     "should return error if login is empty",
			username: "Andrey",
			login:    "",
			password: "test",
			store:    emptyStore,
			code:     codes.InvalidArgument,
		},
		{
			name:     "should return error if password is empty",
			username: "Andrey",
			login:    "admin",
			password: "",
			store:    emptyStore,
			code:     codes.InvalidArgument,
		},
		{
			name:     "should return error if user already exists",
			username: "Andrey",
			login:    "admin",
			password: "test",
			store:    storeWithUser,
			code:     codes.AlreadyExists,
		},
		{
			name:     "should return error if password is too long",
			username: "Andrey",
			login:    "admin",
			password: strings.Repeat("a", 80),
			store:    emptyStore,
			code:     codes.InvalidArgument,
		},
		{
			name:     "should return without error if all data is correct",
			username: "Andrey",
			login:    "admin",
			password: "test",
			store:    emptyStore,
			code:     codes.OK,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := pb.RegisterRequest{
				Username: tc.username,
				Login:    tc.login,
				Password: tc.password,
			}

			server := auth.NewServer(tc.store, jwtManager)

			_, err := server.Register(context.Background(), &req)
			if tc.code == codes.OK {
				require.NoError(t, err)

				user, innerErr := emptyStore.GetUser(tc.login)
				require.NoError(t, innerErr)

				require.Equal(t, user.Username, tc.username)
				require.Equal(t, user.Login, tc.login)

				innerErr = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(tc.password))
				require.NoError(t, innerErr)
			} else {
				require.Error(t, err)
				if e, ok := status.FromError(err); ok {
					require.Equal(t, tc.code, e.Code())
				}
			}
		})
	}
}

func TestServerLogin(t *testing.T) {
	t.Parallel()

	jwtManager := auth.NewJWTManager("secretKey", 15*time.Minute)

	storeWithUser := user.NewInMemoryUserStore()
	newUser, err := domain.NewUser("Andrey", "admin", "test")
	require.NoError(t, err)

	err = storeWithUser.Save(newUser)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		login    string
		password string
		code     codes.Code
	}{
		{
			name:     "should return error if login is empty",
			login:    "",
			password: "test",
			code:     codes.InvalidArgument,
		},
		{
			name:     "should return error if password is empty",
			login:    "admin",
			password: "",
			code:     codes.InvalidArgument,
		},
		{
			name:     "should return error if user not found",
			login:    "admin2",
			password: "test",
			code:     codes.NotFound,
		},
		{
			name:     "should return error if password is incorrect",
			login:    "admin",
			password: "test2",
			code:     codes.Unauthenticated,
		},
		{
			name:     "should return without error if all data is correct",
			login:    "admin",
			password: "test",
			code:     codes.OK,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := pb.LoginRequest{
				Login:    tc.login,
				Password: tc.password,
			}

			server := auth.NewServer(storeWithUser, jwtManager)

			_, err := server.Login(context.Background(), &req)
			if tc.code == codes.OK {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if e, ok := status.FromError(err); ok {
					require.Equal(t, tc.code, e.Code())
				}
			}
		})
	}
}
