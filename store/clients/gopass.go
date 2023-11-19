package clients

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/gopasspw/gopass/pkg/gopass/apimock"
	"github.com/gopasspw/gopass/pkg/gopass/secrets/secparse"
	"github.com/gopasspw/gopass/pkg/otp"
	potp "github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
)

// Gopass Client
type GopassClient interface {
	Get(context.Context, string, string) (gopass.Secret, error)
}

type GopassSecretResolver struct {
	gopass  GopassClient
	context context.Context
}

// ensure that GopassSecretResolver implements the SecretsResolver interface
var _ SecretsResolver = (*GopassSecretResolver)(nil)

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

func (g *GopassSecretResolver) GetOtp(path string, timestamp time.Time) (string, error) {
	secret, err := g.get( /*options*/
		GopassClientGetOptions{
			version: "latest",
			path:    path,
		},
	)
	if err != nil {
		return "", err
	}

	token, err := ResolveOtp(secret, timestamp)
	if err != nil {
		return "", fmt.Errorf("failed to calculate totp token: %s", err)
	}

	return token, nil
}

func NewGopassResolver() GopassSecretResolver {

	ctx := context.Background()
	store, err := api.New(ctx)
	if err != nil {
		fmt.Printf("Failed to initialize gopass API: %s\n", err)
		os.Exit(1)
	}

	return CreateResolver(store, ctx)
}

func CreateResolver(api GopassClient, context context.Context) GopassSecretResolver {
	gopassClient := GopassSecretResolver{}
	gopassClient.context = context
	gopassClient.gopass = api

	return gopassClient
}

func ResolveOtp(secret gopass.Secret, timestamp time.Time) (string, error) {

	token, err := otp.Calculate("_", secret)
	if err != nil {
		return "", fmt.Errorf("failed to calculate totp token")
	}

	code, err := totp.GenerateCodeCustom(
		token.Secret(),
		timestamp,
		totp.ValidateOpts{
			Period:    uint(token.Period()),
			Skew:      1,
			Digits:    potp.DigitsSix,
			Algorithm: potp.AlgorithmSHA1,
		})
	if err != nil {
		return "", fmt.Errorf("failed to generate totp code")
	}

	return code, nil
}

type MockStoredGopassSecret struct {
	Name   []string
	Secret gopass.Secret
}

func MockGopassSecret(t *testing.T, in string) gopass.Secret {
	t.Helper()
	sec, err := secparse.Parse([]byte(in))
	require.NoError(t, err)
	return sec
}

func MockGopassApi(secrets []MockStoredGopassSecret) (GopassSecretResolver, error) {
	ctx := context.Background()
	store := apimock.New()
	for _, sec := range secrets {
		err := store.Set(ctx, strings.Join(sec.Name, "/"), sec.Secret)
		if err != nil {
			return GopassSecretResolver{}, err
		}
	}
	fmt.Print(store.List(ctx))
	return CreateResolver(store, ctx), nil
}
