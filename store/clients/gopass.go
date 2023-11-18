package clients

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/gopasspw/gopass/pkg/gopass/apimock"
	"github.com/gopasspw/gopass/pkg/gopass/secrets"
	"github.com/gopasspw/gopass/pkg/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

// Gopass Client
type GopassClient interface {
	Get(context.Context, string, string) (gopass.Secret, error)
}

type GopassSecretResolver struct {
	gopass  GopassClient
	context context.Context
}

type GopassClientGetOptions struct {
	version string
	path    string
}

func (g *GopassSecretResolver) get(options GopassClientGetOptions) (gopass.Secret, error) {
	var secret gopass.Secret

	version := options.version
	if version == "" {
		version = "latest"
	}
	path := options.path
	if path == "" {
		return secret, errors.New("path is required")
	}

	secret, err := g.gopass.Get(g.context, path, version)

	if err != nil {
		return secret, fmt.Errorf("failed to get secret %s", path)
	}

	return secret, nil
}

func (g *GopassSecretResolver) GetPassword(path string) (string, error) {
	secret, err := g.get( /*options*/
		GopassClientGetOptions{
			version: "latest",
			path:    path,
		},
	)

	if err != nil {
		return "", err
	}

	username := secret.Password()

	return username, nil
}

func (g *GopassSecretResolver) GetUsername(path string) (string, error) {
	secret, err := g.get( /*options*/
		GopassClientGetOptions{
			version: "latest",
			path:    path,
		},
	)

	if err != nil {
		return "", err
	}

	username, exists := secret.Get("username")
	if !exists {
		return "", fmt.Errorf("username not found")
	}

	return username, nil
}

func (g *GopassSecretResolver) GetOtp(path string) (string, error) {
	secret, err := g.get( /*options*/
		GopassClientGetOptions{
			version: "latest",
			path:    path,
		},
	)

	if err != nil {
		return "", err
	}

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

func NewGopassResolver() GopassSecretResolver {
	gopassClient := GopassSecretResolver{}
	gopassClient.context = context.Background()
	api, err := api.New(gopassClient.context)
	if err != nil {
		fmt.Printf("Failed to initialize gopass API: %s\n", err)
		os.Exit(1)
	}
	gopassClient.gopass = api

	return gopassClient
}

func NewMockGopassResolver(data map[string]interface{}) GopassSecretResolver {
	ctx := context.Background()
	gopassClient := GopassSecretResolver{}
	gopassClient.context = ctx
	api := apimock.New()
	for k, v := range data {
		if v == nil {
			continue
		}

		secret := secrets.New()
		// if v.password  set password
		if v.(map[string]interface{})["password"] != nil {
			secret.SetPassword(v.(map[string]interface{})["password"].(string))
		}

		// if v.username  set username
		if v.(map[string]interface{})["username"] != nil {
			secret.Set("username", v.(map[string]interface{})["username"].(string))
		}

		// // if v.totp      set totp
		// if v.(map[string]interface{})["totp"] != nil {
		// 	secret.Set("totp", v.(map[string]interface{})["totp"].(string))
		// }

		api.Set(ctx, k, secret)
	}

	gopassClient.gopass = api

	return gopassClient
}
