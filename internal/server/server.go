package server

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/kaa-it/gophkeeper/internal/server/application/grpc/keeper"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/file"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/kaa-it/gophkeeper/internal/pb"
	"github.com/kaa-it/gophkeeper/internal/server/application/grpc/auth"
	"github.com/kaa-it/gophkeeper/internal/server/infrastructure/storage/user"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

// /gophkeeper.AuthService/Login

func protectedMethods() map[string]bool {
	const keeperServicePath = "/gophkeeper.KeeperService/"

	return map[string]bool{
		keeperServicePath + "UploadCredentials": true,
		keeperServicePath + "UploadFile":        true,
	}
}

// func loadTLSCredentials() (credentials.TransportCredentials, error) {
//	pemClientCA, err := os.ReadFile("cert/ca-cert.pem")
//	if err != nil {
//		return nil, err
//	}
//
//	certPool := x509.NewCertPool()
//	if !certPool.AppendCertsFromPEM(pemClientCA) {
//		return nil, fmt.Errorf("cannot append client CA")
//	}
//
//	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
//	if err != nil {
//		return nil, err
//	}
//
//	config := &tls.Config{
//		MinVersion:   tls.VersionTLS12,
//		Certificates: []tls.Certificate{serverCert},
//		ClientAuth:   tls.RequireAndVerifyClientCert,
//		ClientCAs:    certPool,
//	}
//
//	return credentials.NewTLS(config), nil
// }

func Run() {
	port := flag.Int("port", 0, "the server port")

	flag.Parse()

	log.Printf("start server on port %d", *port)

	userStore := user.NewInMemoryUserStore()
	fileStore := file.NewInMemoryFileStore("files")
	jwtManager := auth.NewJWTManager(secretKey, tokenDuration)

	authServer := auth.NewServer(userStore, jwtManager)
	keeperServer := keeper.NewServer(fileStore)

	// tlsCredentials, err := loadTLSCredentials()
	// if err != nil {
	//	 log.Fatal("cannot load TLS credentials: ", err)
	// }

	interceptor := auth.NewAuthInterceptor(jwtManager, protectedMethods())

	grpcServer := grpc.NewServer(
		// grpc.Creds(tlsCredentials),
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
