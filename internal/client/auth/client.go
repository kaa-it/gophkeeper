package auth

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/kaa-it/gophkeeper/internal/pb"
)

const (
	requestTimeout = 5 * time.Second
)

type Client struct {
	service  pb.AuthServiceClient
	login    string
	password string
}

func NewClient(cc *grpc.ClientConn, login, password string) *Client {
	service := pb.NewAuthServiceClient(cc)

	return &Client{
		service:  service,
		login:    login,
		password: password,
	}
}

func (client *Client) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req := &pb.LoginRequest{
		Login:    client.login,
		Password: client.password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
