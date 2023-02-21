package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChatgptClient(t *testing.T) {
	client := NewChatgptClient("test", "test")
	assert.NotNil(t, client)
}
