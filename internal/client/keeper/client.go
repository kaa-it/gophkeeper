package keeper

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"

	"github.com/kaa-it/gophkeeper/internal/pb"
)

const (
	requestTimeout = 5 * time.Second
	bufferSize     = 1024
)

type Client struct {
	service pb.KeeperServiceClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		service: pb.NewKeeperServiceClient(conn),
	}
}

func (client *Client) UploadCredentials(credentials *pb.Credentials) (string, error) {
	req := &pb.UploadCredentialsRequest{
		Credentials: credentials,
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	res, err := client.service.UploadCredentials(ctx, req)
	if err != nil {
		log.Println("failed to upload credentials: ", err)
		return "", err
	}

	return res.GetId(), nil
}

func (client *Client) UploadFile(filePath string, metadata string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	stream, err := client.service.UploadFile(ctx)
	if err != nil {
		return "", fmt.Errorf("cannot upload file: %w", err)
	}

	req := &pb.UploadFileRequest{
		Data: &pb.UploadFileRequest_Info{
			Info: &pb.FileInfo{
				Metadata: metadata,
				Name:     filepath.Base(filePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		return "", fmt.Errorf("cannot upload file info: %w", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, bufferSize)

	for {
		n, innerErr := reader.Read(buffer)
		if errors.Is(innerErr, io.EOF) {
			break
		}

		if innerErr != nil {
			return "", fmt.Errorf("cannot read chunk to buffer: %w", err)
		}

		req := &pb.UploadFileRequest{
			Data: &pb.UploadFileRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			return "", fmt.Errorf("cannot send chunk to server: %w", err)
		}
	}

	log.Println("close stream")

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("cannot close stream: %w", err)
	}

	return res.GetId(), nil
}
