package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/kaa-it/gophkeeper/internal/pb"
)

const (
	requestTimeout = 5 * time.Second
)

// ErrInvalidArgument is returned when a function receives an invalid argument, such as an empty string or nil value.
var ErrInvalidArgument = errors.New("invalid argument")

// Client manages authentication interactions, including login, registration, and token handling with the AuthService.
type Client struct {
	service       pb.AuthServiceClient
	mutex         sync.Mutex
	login         string
	password      string
	accessToken   string
	notifyChannel chan struct{}
	notifyOnce    sync.Once
}

// NewClient creates a new authentication client for interactions with
// the AuthService using the provided gRPC connection.
func NewClient(cc *grpc.ClientConn) *Client {
	service := pb.NewAuthServiceClient(cc)

	return &Client{
		service:       service,
		notifyChannel: make(chan struct{}),
	}
}

// Login authenticates the client with the given login and password,
// retrieves a new access token, and updates internal state.
func (client *Client) Login(login, password string) error {
	if login == "" || password == "" {
		return fmt.Errorf("empty login or pasword: %w", ErrInvalidArgument)
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.login = login
	client.password = password

	req := &pb.LoginRequest{
		Login:    client.login,
		Password: client.password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	token, err := client.refreshTokenInternal(ctx, req)
	if err != nil {
		return err
	}

	client.accessToken = token

	client.notifyOnce.Do(func() {
		client.notifyChannel <- struct{}{}
	})

	return nil
}

// Register registers a new user with the provided username, login,
// and password and returns an error if registration fails.
func (client *Client) Register(username, login, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req := &pb.RegisterRequest{
		Username: username,
		Login:    login,
		Password: password,
	}

	_, err := client.service.Register(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// AccessToken returns the current access token used for authentication. It is thread-safe.
func (client *Client) AccessToken() string {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	return client.accessToken
}

func (client *Client) refreshToken() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	client.mutex.Lock()

	req := &pb.LoginRequest{
		Login:    client.login,
		Password: client.password,
	}

	client.mutex.Unlock()

	return client.refreshTokenInternal(ctx, req)
}

func (client *Client) refreshTokenInternal(ctx context.Context, req *pb.LoginRequest) (string, error) {
	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
