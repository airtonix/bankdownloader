package clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeychainClient(t *testing.T) {
	store := NewMockKeychainSecretResolver(
		[]MockStoredKeychainSecret{
			{
				Name:     "servicename",
				Username: "someguy",
				Password: "somepassword",
			},
		},
	)
	t.Logf("store: %v", store)
	username, err := store.GetUsername("servicename/someguy")
	assert.NoError(t, err)
	assert.Equal(t, "someguy", username)
	password, err := store.GetPassword("servicename/someguy")
	assert.NoError(t, err)
	assert.Equal(t, "somepassword", password)

}
