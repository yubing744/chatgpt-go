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

// Logger is used for logging formatted messages.
type Logger interface {
	// Printf must have the same semantics as log.Printf.
	Printf(format string, args ...interface{})
}

type ChatgptClient struct {
	session *httpx.HttpSession
	auth    *auth.Authenticator
	logger  Logger
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
		logger:  &Log{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	client := &ChatgptClient{
		baseURL: cfg.baseURL,
		debug:   cfg.debug,
		logger:  cfg.logger,
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
				client.logger.Printf("stop ticker ...\n")
				ticker.Stop()
				return
			case <-ticker.C:
				// 执行刷新 token 的逻辑
				err := client.refreshToken()
				if err != nil {
					client.logger.Printf("fresh token error: %s\n", err.Error())
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
	client.logger.Printf("fresh token ...\n")

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

	client.logger.Printf("fresh token ok!\n")

	return nil
}
