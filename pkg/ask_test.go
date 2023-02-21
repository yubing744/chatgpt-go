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
		Debug:    true,
	}
	client := NewChatgptClient(cfg)
	assert.NotNil(t, client)

	err := client.Start(context.Background())
	defer client.Stop()

	assert.NoError(t, err)

	result, err := client.Ask(context.Background(), "Hello", nil, nil, time.Second*5)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestChatgptClientAsk2(t *testing.T) {
	email := os.Getenv("CHATGPT_EMAIL")
	password := os.Getenv("CHATGPT_PASSWORD")
	assert.NotEmpty(t, email)
	assert.NotEmpty(t, password)

	cfg := &config.Config{
		Email:    email,
		Password: password,
		Proxy:    "",
		Timeout:  time.Second * 30,
		Debug:    true,
	}
	client := NewChatgptClient(cfg)
	assert.NotNil(t, client)

	err := client.Start(context.Background())
	defer client.Stop()
	assert.NoError(t, err)

	result, err := client.Ask(context.Background(), "openAI API 接口 模型温度如何设置？", nil, nil, time.Second*5)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
