package agent

import (
	"context"
	"fmt"
	"log"

	"github.com/dln/luxaforce/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func newLuxaforceClient(hostAndPort, caPath string) (api.LuxaforceClient, error) {
	creds, err := credentials.NewClientTLSFromFile(caPath, "")
	if err != nil {
		return nil, fmt.Errorf("load cert: %v", err)
	}

	conn, err := grpc.Dial(hostAndPort, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("dail: %v", err)
	}
	return api.NewLuxaforceClient(conn), nil
}

// func main() {
// 	client, err := newLuxaforceClient("127.0.0.1:5557", "./grpc.crt")
// 	if err != nil {
// 		log.Fatalf("failed creating luxaforce client: %v ", err)
// 	}
//
// 	req := &api.CreateClientReq{
// 		Client: &api.Client{
// 			Id:     "luxa1",
// 			Name:   "Luxafor Agent 1",
// 			Secret: "fdjkfsdjfskfsdjkfjfdskj",
// 			Labels: []string{"foo=bar"},
// 		},
// 	}
//
// 	if _, err := client.CreateClient(context.TODO(), req); err != nil {
// 		log.Fatalf("failed creating client: %v", err)
// 	}
// }
