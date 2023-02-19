package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yubing744/chatgpt-go/pkg/auth"
)

func TestAuth(t *testing.T) {
	email := os.Getenv("CHATGPT_EMAIL")
	password := os.Getenv("CHATGPT_PASSWORD")
	assert.NotEmpty(t, email)
	assert.NotEmpty(t, password)

	auth := auth.NewAuthenticator(email, password, "")
	assert.NotNil(t, auth)

	err := auth.Begin()
	assert.NoError(t, err)

	accessToken, err := auth.GetAccessToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
}
