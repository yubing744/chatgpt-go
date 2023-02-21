package test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yubing744/chatgpt-go/pkg"
)

func TestChatgptClientLogin(t *testing.T) {
	email := os.Getenv("CHATGPT_EMAIL")
	password := os.Getenv("CHATGPT_PASSWORD")
	assert.NotEmpty(t, email)
	assert.NotEmpty(t, password)

	client := pkg.NewChatgptClient(email, password)
	assert.NotNil(t, client)

	err := client.Start(context.Background())
	defer client.Stop()

	assert.NoError(t, err)
}
