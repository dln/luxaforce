package server

import (
	"context"
	"log"
	"strings"

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

type ValidateToken struct {
	GoogleClientID string
}

func NewValidateToken(clientId string) *ValidateToken {
	return &ValidateToken{
		GoogleClientID: clientId,
	}
}

func (v *ValidateToken) valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	verifier := googleAuthIDTokenVerifier.Verifier{}

	err := verifier.VerifyIDToken(token, []string{
		v.GoogleClientID,
	})
	if err != nil {
		log.Printf("Token verified failed: %v", err)
		return false
	}
	// claimSet.Iss,claimSet.Email ...
	claimSet, err := googleAuthIDTokenVerifier.Decode(token)
	if err != nil {
		return false
	}
	log.Printf("User authenicated: %s", claimSet.Email)
	return true
}

func (v *ValidateToken) EnsureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !v.valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
