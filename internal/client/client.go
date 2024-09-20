package client

import (
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kaa-it/gophkeeper/internal/client/auth"
	"github.com/kaa-it/gophkeeper/internal/client/keeper"
)

// var ErrCannotAppendServerCA = errors.New("cannot append server CA")

const (
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const keeperServicePath = "/gophkeeper.KeeperService/"

	return map[string]bool{
		keeperServicePath + "UploadCredentials": true,
		keeperServicePath + "UploadFile":        true,
	}
}

// func loadTLSCredentials() (credentials.TransportCredentials, error) {
//	pemServerCA, err := os.ReadFile("cert/ca-cert.pem")
//	if err != nil {
//		return nil, err
//	}
//
//	certPool := x509.NewCertPool()
//	if !certPool.AppendCertsFromPEM(pemServerCA) {
//		return nil, ErrCannotAppendServerCA
//	}
//
//	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
//	if err != nil {
//		return nil, err
//	}
//
//	config := &tls.Config{
//		MinVersion:   tls.VersionTLS12,
//		Certificates: []tls.Certificate{clientCert},
//		RootCAs:      certPool,
//	}
//
//	return credentials.NewTLS(config), nil
// }

// func testRateLaptop(laptopClient *client.LaptopClient) {
//	n := 3
//	laptopIDs := make([]string, n)
//	scores := make([]float64, n)
//
//	for i := 0; i < n; i++ {
//		laptop := sample.NewLaptop()
//		laptopClient.CreateLaptop(laptop)
//		laptopIDs[i] = laptop.GetId()
//	}
//
//	for {
//		fmt.Print("rate laptop (y/n)? ")
//		var answer string
//		fmt.Scan(&answer)
//
//		if strings.ToLower(answer) != "y" {
//			break
//		}
//
//		for i := 0; i < n; i++ {
//			scores[i] = sample.RandomLaptopScore()
//		}
//
//		err := laptopClient.RateLaptop(laptopIDs, scores)
//		if err != nil {
//			log.Fatal("cannot rate laptop: ", err)
//		}
//	}
//}

func Run() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()

	log.Printf("dial server %s", *serverAddress)

	// tlsCredentials, err := loadTLSCredentials()
	// if err != nil {
	//	log.Fatal("cannot load TLS credentials: ", err)
	// }
	//
	// cc1, err := grpc.NewClient(*serverAddress, grpc.WithTransportCredentials(tlsCredentials))
	// if err != nil {
	//	log.Fatal("cannot dial server: ", err)
	// }
	// defer cc1.Close()

	cc1, err := grpc.NewClient(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}

	authClient := auth.NewClient(cc1)
	interceptor := auth.NewAuthInterceptor(authClient, authMethods(), refreshDuration)

	cc2, err := grpc.NewClient(
		*serverAddress,
		// grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	defer cc2.Close()

	keeperClient := keeper.NewClient(cc2)

	if regErr := authClient.Register("admin", "admin", "admin"); regErr != nil {
		log.Printf("cannot register: %v", err)
		return
	}

	authClient.Login("admin", "admin")
	// keeperClient.UploadCredentials(&pb.Credentials{
	//	Metadata: "Yandex",
	//	Login:    "Test",
	//	Password: "Test",
	// })
	fileID, err := keeperClient.UploadFile("tmp/laptop.jpg", "Laptop")
	if err != nil {
		log.Printf("cannot upload file: %v", err)
		return
	}

	log.Printf("upload file to %s", fileID)
}
