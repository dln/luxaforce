package server

import (
	"context"
	"errors"
	"strings"

	"github.com/dln/luxaforce/api"
	"github.com/futurenda/google-auth-id-token-verifier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

type luxaforceAPI struct{}

func NewAPI() api.LuxaforceServer {
	return luxaforceAPI{}
}

func (d luxaforceAPI) CreateClient(ctx context.Context, req *api.CreateClientReq) (*api.CreateClientResp, error) {
	if req.Client == nil {
		return nil, errors.New("no client supplied")
	}
	return nil, nil
}
func (d luxaforceAPI) DeleteClient(ctx context.Context, req *api.DeleteClientReq) (*api.DeleteClientResp, error) {
	return nil, nil
}

func (d luxaforceAPI) GetVersion(ctx context.Context, req *api.VersionReq) (*api.VersionResp, error) {
	return &api.VersionResp{
		Server: "0.0.1",
		Api:    1,
	}, nil
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	v := googleAuthIDTokenVerifier.Verifier{}
	aud := "xxxxxx-yyyyyyy.apps.googleusercontent.com"
	err := v.VerifyIDToken(token, []string{
		aud,
	})
	if err != nil {
		return false
	}
	// claimSet.Iss,claimSet.Email ...
	_, err = googleAuthIDTokenVerifier.Decode(token)
	if err != nil {
		return false
	}

	return true
}

func EnsureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
