package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dln/luxaforce/agent"
	"github.com/dln/luxaforce/api"
	"github.com/dln/luxaforce/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app                = kingpin.New("luxaforce", "Control the light.")
	googleClientId     = app.Flag("google-client-id", "Google Client ID.").Required().String()
	googleClientSecret = app.Flag("google-client-secret", "Google Secret.").Required().String()
	srv                = app.Command("server", "Start a new server.")
	agnt               = app.Command("agent", "Start a agent.")
)

func main() {
	grpcAddr := "127.0.0.1:8888"

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case srv.FullCommand():
		cert, err := tls.LoadX509KeyPair(testdata.Path("server1.pem"), testdata.Path("server1.key"))
		if err != nil {
			log.Fatalf("failed to load key pair: %s", err)
		}
		validateToken := server.NewValidateToken(*googleClientId)
		grpcOptions := []grpc.ServerOption{
			// Validate Google OAuth token
			grpc.UnaryInterceptor(validateToken.EnsureValidToken),
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
	case agnt.FullCommand():
		client, err := agent.NewLuxaforceClient(grpcAddr, *googleClientId, *googleClientSecret)
		if err != nil {
			log.Fatalf("Client failed: %v", err)
		}
		req := &api.CreateClientReq{
			Client: &api.Client{
				Id:     "luxa1",
				Name:   "Luxafor Agent 1",
				Secret: "",
				Labels: []string{"foo=bar"},
			},
		}

		if _, err := client.CreateClient(context.TODO(), req); err != nil {
			log.Fatalf("failed creating client: %v", err)
		}
	}
}
