package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/creditcard"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/text"

	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/credentials"

	"github.com/kaa-it/gophkeeper/internal/server/application/grpc/keeper"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/file"

	"google.golang.org/grpc"
	grpcCredentials "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/kaa-it/gophkeeper/internal/pb"
	"github.com/kaa-it/gophkeeper/internal/server/application/grpc/auth"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/user"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

// ErrCannotAppendClientCA is the error returned when the client CA certificate cannot be appended to the cert pool.
var ErrCannotAppendClientCA = errors.New("cannot append client CA")

// ProtectedMethods returns a map of gRPC service methods
// that are protected and require authentication.
func ProtectedMethods() map[string]bool {
	const keeperServicePath = "/gophkeeper.KeeperService/"

	return map[string]bool{
		keeperServicePath + "UploadCredentials": true,
		keeperServicePath + "UploadFile":        true,
	}
}

func loadTLSCredentials() (grpcCredentials.TransportCredentials, error) {
	pemClientCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, ErrCannotAppendClientCA
	}

	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return grpcCredentials.NewTLS(config), nil
}

func Run() {
	port := flag.Int("port", 0, "the server port")

	flag.Parse()

	log.Printf("start server on port %d", *port)

	userStore := user.NewInMemoryUserStore()
	fileStore := file.NewInMemoryFileStore("files")
	credentialsStore := credentials.NewInMemoryCredentialsStore()
	creditCardStore := creditcard.NewInMemoryCreditCardStore()
	textStore := text.NewInMemoryTextStore()
	jwtManager := auth.NewJWTManager(secretKey, tokenDuration)

	authServer := auth.NewServer(userStore, jwtManager)
	keeperServer := keeper.NewServer(fileStore, credentialsStore, creditCardStore, textStore)

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	interceptor := auth.NewAuthInterceptor(jwtManager, ProtectedMethods())

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterKeeperServiceServer(grpcServer, keeperServer)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
