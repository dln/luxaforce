package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/dln/luxaforce/api"
	"github.com/dln/luxaforce/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

func main() {
	grpcAddr := "127.0.0.1:5557"
	cert, err := tls.LoadX509KeyPair(testdata.Path("server1.pem"), testdata.Path("server1.key"))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	grpcOptions := []grpc.ServerOption{
		// Validate Google OAuth token
		grpc.UnaryInterceptor(server.EnsureValidToken),
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	errc := make(chan error, 1)
	log.Printf("listening (grpc) on %s", grpcAddr)
	go func() {
		errc <- func() error {
			list, err := net.Listen("tcp", grpcAddr)
			if err != nil {
				return fmt.Errorf("listening on %s failed: %v", grpcAddr, err)
			}
			s := grpc.NewServer(grpcOptions...)
			api.RegisterLuxaforceServer(s, server.NewAPI())
			err = s.Serve(list)
			return fmt.Errorf("listening on %s failed: %v", grpcAddr, err)
		}()
	}()
	<-errc
}
