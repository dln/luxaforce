package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dln/luxaforce/api"
	"github.com/dln/luxaforce/server"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
)

func main() {
	grpcAddr := "127.0.0.1:5557"
	var grpcOptions []grpc.ServerOption

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
