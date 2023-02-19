package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthenticator(t *testing.T) {
	email := os.Getenv("CHATGPT_EMAIL")
	password := os.Getenv("CHATGPT_PASSWORD")
	assert.NotEmpty(t, email)
	assert.NotEmpty(t, password)

	auth := NewAuthenticator(email, password, "")
	assert.NotNil(t, auth)
}
