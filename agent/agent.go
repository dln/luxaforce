package agent

import (
	"crypto/tls"
	"fmt"

	"github.com/dln/luxaforce/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func NewLuxaforceClient(hostAndPort string) (api.LuxaforceClient, error) {
	login := &LoginAgent{
		AllowBrowser: true,
	}

	token, err := login.PerformLogin()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(hostAndPort, []grpc.DialOption{
		grpc.WithPerRPCCredentials(oauth.NewOauthAccess(token)),
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		),
	}...)
	if err != nil {
		return nil, fmt.Errorf("dail: %v", err)
	}

	return api.NewLuxaforceClient(conn), nil
}

// func main() {
// 	client, err := newLuxaforceClient("127.0.0.1:5557")
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
