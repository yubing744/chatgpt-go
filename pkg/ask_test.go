package pkg

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yubing744/chatgpt-go/pkg/config"
)

func TestChatgptClientAsk(t *testing.T) {
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

	result, err := client.Ask(context.Background(), "Hello", nil, nil, time.Second*5)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
