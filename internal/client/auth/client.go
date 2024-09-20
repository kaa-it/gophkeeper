package auth

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/kaa-it/gophkeeper/internal/pb"
)

const (
	requestTimeout = 5 * time.Second
)

type Client struct {
	service       pb.AuthServiceClient
	mutex         sync.Mutex
	login         string
	password      string
	notifyChannel chan struct{}
}

func NewClient(cc *grpc.ClientConn) *Client {
	service := pb.NewAuthServiceClient(cc)

	return &Client{
		service:       service,
		notifyChannel: make(chan struct{}, 1),
	}
}

func (client *Client) Login(login, password string) {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.login = login
	client.password = password

	client.notifyChannel <- struct{}{}
}

func (client *Client) RefreshToken() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	client.mutex.Lock()

	req := &pb.LoginRequest{
		Login:    client.login,
		Password: client.password,
	}

	client.mutex.Unlock()

	if len(client.login) == 0 || len(client.password) == 0 {
		return "", nil
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}

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
