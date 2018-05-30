package agent

import (
	"crypto/tls"
	"fmt"

	"github.com/dln/luxaforce/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewLuxaforceClient(hostAndPort string) (api.LuxaforceClient, error) {
	login := &LoginAgent{
		AllowBrowser: true,
		ClientID:     "x-y.apps.googleusercontent.com",
		ClientSecret: "notsosecret",
	}

	ts, err := login.PerformLogin()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(hostAndPort, []grpc.DialOption{
		grpc.WithPerRPCCredentials(NewOauthAccess(ts)),
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		),
	}...)
	if err != nil {
		return nil, fmt.Errorf("dail: %v", err)
	}

	return api.NewLuxaforceClient(conn), nil
}
