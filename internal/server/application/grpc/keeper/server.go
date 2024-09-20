package keeper

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kaa-it/gophkeeper/internal/pb"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/file"
)

const maxFileSize = 1024 * 1024 * 10

type Server struct {
	fileStore file.Repository
	pb.UnimplementedKeeperServiceServer
}

func NewServer(fileStore file.Repository) *Server {
	return &Server{fileStore: fileStore}
}

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
