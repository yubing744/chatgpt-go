package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yubing744/chatgpt-go/pkg/config"
)

func TestNewChatgptClient(t *testing.T) {
	email := os.Getenv("CHATGPT_EMAIL")
	password := os.Getenv("CHATGPT_PASSWORD")
	assert.NotEmpty(t, email)
	assert.NotEmpty(t, password)

	cfg := &config.Config{
		Email:    email,
		Password: password,
		Proxy:    "",
	}
	client := NewChatgptClient(cfg)
	assert.NotNil(t, client)
}

func TestChatgptClientLogin(t *testing.T) {
	email := os.Getenv("CHATGPT_EMAIL")
	password := os.Getenv("CHATGPT_PASSWORD")
	assert.NotEmpty(t, email)
	assert.NotEmpty(t, password)

	cfg := &config.Config{
		Email:    email,
		Password: password,
		Proxy:    "",
	}
	client := NewChatgptClient(cfg)
	assert.NotNil(t, client)

	err := client.Login()
	assert.NoError(t, err)
}
