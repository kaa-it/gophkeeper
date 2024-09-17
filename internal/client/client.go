package client

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/kaa-it/gophkeeper/internal/client/auth"
)

const (
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const laptopServicePath = "/LaptopService/"

	return map[string]bool{
		laptopServicePath + "CreateLaptop": true,
		laptopServicePath + "UploadLaptop": true,
		laptopServicePath + "RateLaptop":   true,
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("cannot append server CA")
	}

	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

//func testRateLaptop(laptopClient *client.LaptopClient) {
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
	//interceptor, err := auth.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	//if err != nil {
	//	log.Fatal("cannot create auth interceptor: ", err)
	//}

	//cc2, err := grpc.NewClient(
	//	*serverAddress,
	//	grpc.WithTransportCredentials(tlsCredentials),
	//	grpc.WithUnaryInterceptor(interceptor.Unary()),
	//	grpc.WithStreamInterceptor(interceptor.Stream()),
	//)
	//if err != nil {
	//	log.Fatal("cannot dial server: ", err)
	//}
	//defer cc2.Close()
	//
	//laptopClient := client.NewLaptopClient(cc2)
}
