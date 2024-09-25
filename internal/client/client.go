package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/kaa-it/gophkeeper/internal/pb"

	"google.golang.org/grpc"

	"github.com/kaa-it/gophkeeper/internal/client/auth"
	"github.com/kaa-it/gophkeeper/internal/client/keeper"
)

// ErrCannotAppendServerCA is the error returned when the server CA certificate cannot be appended to the cert pool.
var ErrCannotAppendServerCA = errors.New("cannot append server CA")

const (
	refreshDuration = 30 * time.Second
)

// AuthMethods returns a map of gRPC method names to a boolean indicating if the method requires authentication.
func AuthMethods() map[string]bool {
	const keeperServicePath = "/gophkeeper.KeeperService/"

	return map[string]bool{
		keeperServicePath + "UploadCredentials": true,
		keeperServicePath + "UploadFile":        true,
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, ErrCannotAppendServerCA
	}

	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

// Run initializes the gRPC clients, registers and logs in the admin user, and attempts to upload a file to the server.
func Run() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()

	log.Printf("dial server %s", *serverAddress)

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	cc1, err := grpc.NewClient(*serverAddress, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	defer cc1.Close()

	authClient := auth.NewClient(cc1)
	interceptor := auth.NewAuthInterceptor(authClient, AuthMethods(), refreshDuration)

	cc2, err := grpc.NewClient(
		*serverAddress,
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Println("cannot dial server: ", err)
		return
	}
	defer cc2.Close()

	keeperClient := keeper.NewClient(cc2)

	if regErr := authClient.Register("admin", "admin", "admin"); regErr != nil {
		log.Printf("cannot register: %v", regErr)
		return
	}

	if logErr := authClient.Login("admin", "admin"); logErr != nil {
		log.Printf("cannot login: %v", logErr)
		return
	}

	credentialsID, err := keeperClient.UploadCredentials(&pb.Credentials{
		Metadata: "Yandex",
		Login:    "Test",
		Password: "Test",
	})
	if err != nil {
		log.Printf("cannot upload credentials: %v", err)
		return
	}

	log.Printf("upload credentials to %s", credentialsID)

	fileID, err := keeperClient.UploadFile("tmp/laptop.jpg", "Laptop")
	if err != nil {
		log.Printf("cannot upload file: %v", err)
		return
	}

	log.Printf("upload file to %s", fileID)
}
