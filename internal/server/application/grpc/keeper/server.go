package keeper

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/creditcard"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/text"

	"github.com/kaa-it/gophkeeper/internal/server/domain"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/credentials"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kaa-it/gophkeeper/internal/pb"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/file"
)

const maxFileSize = 1024 * 1024 * 10

// Server handles incoming requests and manages file storage operations.
type Server struct {
	fileStore        file.Repository
	credentialsStore credentials.Repository
	creditCardStore  creditcard.Repository
	textStore        text.Repository
	pb.UnimplementedKeeperServiceServer
}

// NewServer initializes and returns a new instance of Server with the provided repositories
// for files, credentials, credit cards, and texts.
func NewServer(
	fileStore file.Repository,
	credentialsStore credentials.Repository,
	creditCardStore creditcard.Repository,
	textStore text.Repository,
) *Server {
	return &Server{
		fileStore:        fileStore,
		credentialsStore: credentialsStore,
		creditCardStore:  creditCardStore,
		textStore:        textStore,
	}
}

// UploadFile handles the client streaming request for uploading a file and stores it in the server, returning an ID.
func (server *Server) UploadFile(stream grpc.ClientStreamingServer[pb.UploadFileRequest, pb.UploadFileResponse]) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info: %v", err))
	}

	metadata := req.GetInfo().GetMetadata()
	fileName := req.GetInfo().GetName()

	log.Printf("receive an upload-file request with metadata: %s, file name: %s", metadata, fileName)

	fileData := bytes.Buffer{}
	fileSize := 0

	for {
		if innerErr := contextError(stream.Context()); innerErr != nil {
			return innerErr
		}

		log.Print("waiting for more file data...")

		req, innerErr := stream.Recv()
		if errors.Is(innerErr, io.EOF) {
			log.Print("no more file data")
			break
		}

		if innerErr != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		fileSize += size

		if fileSize > maxFileSize {
			return logError(status.Errorf(codes.InvalidArgument, "file size is too large: %d", fileSize))
		}

		_, err = fileData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot write chunk data: %v", err))
		}
	}

	fileID, err := server.fileStore.Save(metadata, fileName, fileData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save file: %v", err))
	}

	res := &pb.UploadFileResponse{
		Id: fileID,
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	log.Printf("saved file with id: %s", fileID)

	return nil
}

// UploadCredentials processes incoming credentials request and
// saves the credentials into the repository, returning an ID.
func (server *Server) UploadCredentials(
	_ context.Context,
	req *pb.UploadCredentialsRequest,
) (*pb.UploadCredentialsResponse, error) {
	requestCredentials := req.GetCredentials()

	newCredentials := domain.NewCredentials(
		requestCredentials.GetLogin(),
		requestCredentials.GetPassword(),
		requestCredentials.GetMetadata(),
	)

	id, err := server.credentialsStore.Save(newCredentials)
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "failed to save credentials: %v", err)
	}

	res := &pb.UploadCredentialsResponse{
		Id: id,
	}

	return res, nil
}

// UploadText handles the request to upload a text document and saves it in the repository, returning an ID.
func (server *Server) UploadText(
	_ context.Context,
	req *pb.UploadTextRequest,
) (*pb.UploadTextResponse, error) {
	requestText := req.GetText()

	newText := domain.NewText(
		requestText.GetMetadata(),
		requestText.GetText(),
	)

	id, err := server.textStore.Save(newText)
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "failed to save text: %v", err)
	}

	res := &pb.UploadTextResponse{
		Id: id,
	}

	return res, nil
}

// UploadCreditCard processes an incoming credit card request,
// saves it to the repository, and returns an ID or an error.
func (server *Server) UploadCreditCard(
	_ context.Context,
	req *pb.UploadCreditCardRequest,
) (*pb.UploadCreditCardResponse, error) {
	requestCreditCard := req.GetCreditCard()

	requestBillingAddress := requestCreditCard.GetBillingAddress()

	billingAddress := domain.NewBillingAddress(
		requestBillingAddress.GetCountry(),
		requestBillingAddress.GetStreet(),
		requestBillingAddress.GetCity(),
		requestBillingAddress.GetPostcode(),
	)

	newCreditCard := domain.NewCreditCard(
		requestCreditCard.GetMetadata(),
		requestCreditCard.GetName(),
		requestCreditCard.GetMonth(),
		requestCreditCard.GetYear(),
		requestCreditCard.GetNumber(),
		*billingAddress,
	)

	id, err := server.creditCardStore.Save(newCreditCard)
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "failed to save credit card: %v", err)
	}

	res := &pb.UploadCreditCardResponse{
		Id: id,
	}

	return res, nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}

	return err
}

func contextError(ctx context.Context) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		return logError(status.Error(codes.Canceled, "request is canceled"))
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	}

	return nil
}
