package clients

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/apimock"
	"github.com/gopasspw/gopass/pkg/gopass/secrets/secparse"
	"github.com/stretchr/testify/require"
)

type MockStoredGopassSecret struct {
	Name   []string
	Secret gopass.Secret
}

func NewMockGopassSecret(t *testing.T, in string) gopass.Secret {
	t.Helper()
	sec, err := secparse.Parse([]byte(in))
	require.NoError(t, err)
	return sec
}

func NewMockGopassSecretResolver(secrets []MockStoredGopassSecret) (*GopassSecretResolver, error) {
	ctx := context.Background()
	store := apimock.New()
	for _, sec := range secrets {
		err := store.Set(ctx, strings.Join(sec.Name, "/"), sec.Secret)
		if err != nil {
			return &GopassSecretResolver{}, err
		}
	}
	fmt.Print(store.List(ctx))
	return CreateResolver(store, ctx), nil
}
