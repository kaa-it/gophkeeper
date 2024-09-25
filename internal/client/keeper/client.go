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

// Client represents a client to interact with KeeperService.
type Client struct {
	service pb.KeeperServiceClient
}

// NewClient creates a new instance of Client using the provided grpc.ClientConn.
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		service: pb.NewKeeperServiceClient(conn),
	}
}

// UploadCredentials uploads the provided credentials to the KeeperService
// and returns an ID or an error if the upload fails.
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

// UploadCreditCard uploads the provided credit card details to the KeeperService
// and returns an ID or an error if the upload fails.
func (client *Client) UploadCreditCard(creditCard *pb.CreditCard) (string, error) {
	req := &pb.UploadCreditCardRequest{
		CreditCard: creditCard,
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	res, err := client.service.UploadCreditCard(ctx, req)
	if err != nil {
		log.Println("failed to upload credit card: ", err)
		return "", err
	}

	return res.GetId(), nil
}

// UploadText uploads the provided text to the KeeperService and returns an ID or an error if the upload fails.
func (client *Client) UploadText(text *pb.Text) (string, error) {
	req := &pb.UploadTextRequest{
		Text: text,
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	res, err := client.service.UploadText(ctx, req)
	if err != nil {
		log.Println("failed to upload text: ", err)
		return "", err
	}

	return res.GetId(), nil
}

// UploadFile uploads a file specified by filePath to the server with the provided metadata.
// It returns the ID of the uploaded file or an error if the upload process fails.
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
