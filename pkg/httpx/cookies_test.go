package httpx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCookiesGet(t *testing.T) {
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

	val, ok := cookies.Get("SUID")
	assert.True(t, ok)
	assert.NotEmpty(t, val)
}

func TestCookiesSet(t *testing.T) {
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

	ok := cookies.Set("SUID", "xxxx")
	assert.True(t, ok)

	val, ok := cookies.Get("SUID")
	assert.True(t, ok)
	assert.Equal(t, "xxxx", val)
}
