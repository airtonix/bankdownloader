package clients

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/gopasspw/gopass/pkg/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

// Gopass Client

type GopassClient struct {
	gopass  gopass.Store
	context context.Context
}
type GopassClientGetOptions struct {
	Version string
	Path    string
	Key     string
}

func (g *GopassClient) Get(options GopassClientGetOptions) (gopass.Secret, error) {
	var secret gopass.Secret

	version := options.Version
	if version == "" {
		version = "latest"
	}
	path := options.Path
	if path == "" {
		return secret, errors.New("path is required")
	}

	secret, err := g.gopass.Get(g.context, path, version)

	if err != nil {
		return secret, fmt.Errorf("failed to get secret %s", path)
	}

	return secret, nil
}

func (g *GopassClient) GetOtpToken(secret gopass.Secret) (string, error) {
	token, err := otp.Calculate("", secret)

	if err != nil {
		return "", fmt.Errorf("failed to calculate totp token")
	}
	switch token.Type() {
	case "totp":
		return totp.GenerateCode(token.Secret(), time.Now())

	case "hotp":
		return hotp.GenerateCode(token.Secret(), token.Period())
	}

	return "", fmt.Errorf("failed to calculate totp token")
}

func NewGopassClient() GopassClient {
	gopassClient := GopassClient{}
	gopassClient.context = context.Background()
	api, err := api.New(gopassClient.context)
	if err != nil {
		fmt.Printf("Failed to initialize gopass API: %s\n", err)
		os.Exit(1)
	}
	gopassClient.gopass = api

	return gopassClient
}
