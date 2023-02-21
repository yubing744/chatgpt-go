package pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/yubing744/chatgpt-go/pkg/auth"
	"github.com/yubing744/chatgpt-go/pkg/httpx"
)

type ChatgptClient struct {
	session *httpx.HttpSession
	auth    *auth.Authenticator
	cancel  context.CancelFunc
	baseURL string
	debug   bool
}

func NewChatgptClient(email string, password string, opts ...Option) *ChatgptClient {
	cfg := &Options{
		baseURL: "https://chatgpt.duti.tech",
		timeout: time.Second * 300,
		proxy:   "",
		debug:   false,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	client := &ChatgptClient{
		baseURL: cfg.baseURL,
		debug:   cfg.debug,
	}

	session, err := httpx.NewHttpSession(cfg.timeout)
	if err != nil {
		log.Fatal("init http session fatal")
	}

	client.session = session
	client.auth = auth.NewAuthenticator(email, password, cfg.proxy)

	return client
}

func (client *ChatgptClient) Start(ctx context.Context) error {
	ctx, client.cancel = context.WithCancel(ctx)

	err := client.auth.Begin()
	if err != nil {
		return errors.Wrap(err, "Error in auth")
	}

	err = client.refreshToken()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(10 * time.Minute) // 每 10 分钟刷新一次 token

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("stop ticker ...\n")
				ticker.Stop()
				return
			case <-ticker.C:
				// 执行刷新 token 的逻辑
				err := client.refreshToken()
				if err != nil {
					fmt.Printf("fresh token error: %s\n", err.Error())
					continue
				}
			}
		}
	}()

	return nil
}

func (client *ChatgptClient) Stop() {
	if client.cancel != nil {
		client.cancel()
	}
}

func (client *ChatgptClient) refreshToken() error {
	fmt.Printf("fresh token ...\n")

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

	fmt.Printf("fresh token ok!\n")

	return nil
}
