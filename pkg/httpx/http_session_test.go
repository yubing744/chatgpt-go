package httpx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpSession(t *testing.T) {
	client, err := NewHttpSession(time.Second * 5)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestHTTPXGet(t *testing.T) {
	client, err := NewHttpSession(time.Second * 5)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	resp, err := client.Get("https://www.bing.com/", nil, true)
	if resp != nil {
		defer resp.Body.Close()
	}

	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
}

func TestHTTPXGetCookies(t *testing.T) {
	client, err := NewHttpSession(time.Second * 5)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	resp, err := client.Get("https://www.bing.com/", nil, true)
	if resp != nil {
		defer resp.Body.Close()
	}

	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	cookies := client.Cookies("bing.com")
	assert.NotNil(t, cookies)
}
