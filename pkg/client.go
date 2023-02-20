package pkg

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/yubing744/chatgpt-go/pkg/auth"
	"github.com/yubing744/chatgpt-go/pkg/config"
	"github.com/yubing744/chatgpt-go/pkg/httpx"
)

type ChatgptClient struct {
	baseURL string
	session *httpx.HttpSession
	auth    *auth.Authenticator
	debug   bool
}

func NewChatgptClient(cfg *config.Config) *ChatgptClient {
	client := &ChatgptClient{
		baseURL: "https://chatgpt.duti.tech",
		debug:   cfg.Debug,
	}

	session, err := httpx.NewHttpSession(cfg.Timeout)
	if err != nil {
		log.Fatal("init http session fatal")
	}

	client.session = session
	client.auth = auth.NewAuthenticator(cfg.Email, cfg.Password, cfg.Proxy)

	return client
}

func (client *ChatgptClient) Login() error {
	err := client.auth.Begin()
	if err != nil {
		return errors.Wrap(err, "Error in auth")
	}

	client.refreshToken()

	return nil
}

func (client *ChatgptClient) refreshToken() error {
	accessToken, err := client.auth.GetAccessToken()
	if err != nil {
		return errors.Wrap(err, "Error in get access token")
	}

	client.session.SetHeaders(http.Header{
		"Accept":                    {"text/event-stream"},
		"Authorization":             {fmt.Sprintf("Bearer %s", accessToken)},
		"Content-Type":              {"application/json"},
		"X-Openai-Assistant-App-Id": {""},
		"Connection":                {"close"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Referer":                   {"https://chat.openai.com/chat"},
	})

	return nil
}
