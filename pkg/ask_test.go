package pkg

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yubing744/chatgpt-go/pkg/config"
)

func TestParseResponse(t *testing.T) {
	cfg := &config.Config{
		Email:    "email",
		Password: "password",
		Proxy:    "",
		Timeout:  time.Second * 30,
		Debug:    true,
	}
	client := NewChatgptClient(cfg)

	body := ``
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	msgs, err := client.parseResponse(resp)
	assert.NoError(t, err)
	assert.NotNil(t, msgs)
}

func TestParseResponseForUnmarshalError(t *testing.T) {
	cfg := &config.Config{
		Email:    "email",
		Password: "password",
		Proxy:    "",
		Timeout:  time.Second * 30,
		Debug:    true,
	}
	client := NewChatgptClient(cfg)

	body := `data: 2023-02-21 07:00:21.653311`
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	msgs, err := client.parseResponse(resp)
	assert.NoError(t, err)
	assert.NotNil(t, msgs)
}

func TestParseResponseForDetail(t *testing.T) {
	cfg := &config.Config{
		Email:    "email",
		Password: "password",
		Proxy:    "",
		Timeout:  time.Second * 30,
		Debug:    true,
	}
	client := NewChatgptClient(cfg)

	body := `{"detail":"Too many requests in 1 hour. Try again later."}`
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	_, err := client.parseResponse(resp)
	assert.Error(t, err)
	assert.Equal(t, "Too many requests in 1 hour. Try again later.", err.Error())
}

func TestParseResponseForServerError(t *testing.T) {
	cfg := &config.Config{
		Email:    "email",
		Password: "password",
		Proxy:    "",
		Timeout:  time.Second * 30,
		Debug:    true,
	}
	client := NewChatgptClient(cfg)

	body := `{"detail":{"message":"The server had an error while processing your request. Sorry about that! You can retry your request, or contact us through our help center at help.openai.com if the error persists. (Please include the request ID 985e0eeb2c44145e93637d2d79d416cf in your message.)","type":"server_error","param":null,"code":null}}`
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	_, err := client.parseResponse(resp)
	assert.Error(t, err)
	assert.Equal(t, `{"detail":{"message":"The server had an error while processing your request. Sorry about that! You can retry your request, or contact us through our help center at help.openai.com if the error persists. (Please include the request ID 985e0eeb2c44145e93637d2d79d416cf in your message.)","type":"server_error","param":null,"code":null}}`, err.Error())
}

func TestParseResponseForInternalServerError(t *testing.T) {
	cfg := &config.Config{
		Email:    "email",
		Password: "password",
		Proxy:    "",
		Timeout:  time.Second * 30,
		Debug:    true,
	}
	client := NewChatgptClient(cfg)

	body := `Internal Server Error`
	resp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}

	_, err := client.parseResponse(resp)
	assert.Error(t, err)
	assert.Equal(t, "Internal Server Error", err.Error())
}
