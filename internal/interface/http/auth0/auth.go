package auth0

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func New() (*Authenticator, error) {
	provier, err := oidc.NewProvider(context.Background(), "https://"+os.Getenv("AUTH0_DOMAIN")+"/")
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provier.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{provier, conf}, nil
}

func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
