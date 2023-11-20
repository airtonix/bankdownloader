package clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeychainClient(t *testing.T) {
	store := NewMockKeychainSecretResolver(
		"test",
		[]MockStoredKeychainSecret{
			{
				Name:   "pathtosecret",
				Secret: "somepassword",
			},
		},
	)
	t.Logf("store: %v", store)
	username, err := store.GetUsername("pathtosecret")
	assert.NoError(t, err)
	assert.Equal(t, "pathtosecret", username)
	password, err := store.GetPassword("pathtosecret")
	assert.NoError(t, err)
	assert.Equal(t, "somepassword", password)

}
