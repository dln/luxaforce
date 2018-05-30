package agent

import (
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
)

type tokenSource struct {
	oauth2.TokenSource
}

func NewOauthAccess(src oauth2.TokenSource) credentials.PerRPCCredentials {
	return tokenSource{src}
}

func (ts tokenSource) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	log.Printf("TOKEN %+v", token)
	return map[string]string{
		"authorization": token.Type() + " " + token.Extra("id_token").(string),
	}, nil
}

func (ts tokenSource) RequireTransportSecurity() bool {
	return true
}
