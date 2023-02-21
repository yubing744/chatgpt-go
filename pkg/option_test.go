package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithOptions(t *testing.T) {
	cfg := &Options{
		baseURL: "https://chatgpt.duti.tech",
	}

	opt := WithOptions(Options{
		baseURL: "https://chatgpt.duti.tech2",
	})

	opt(cfg)

	assert.Equal(t, "https://chatgpt.duti.tech2", cfg.baseURL)
}

func TestWithBaseURL(t *testing.T) {
	cfg := &Options{
		baseURL: "https://chatgpt.duti.tech",
	}

	opt := WithBaseURL("https://chatgpt.duti.tech2")

	opt(cfg)

	assert.Equal(t, "https://chatgpt.duti.tech2", cfg.baseURL)
}

func TestWithProxy(t *testing.T) {
	cfg := &Options{
		proxy: "",
	}

	opt := WithProxy("127.0.0.1:8081")

	opt(cfg)

	assert.Equal(t, "127.0.0.1:8081", cfg.proxy)
}

func TestWithTimeout(t *testing.T) {
	cfg := &Options{
		timeout: time.Second * 5,
	}

	opt := WithTimeout(time.Second * 10)

	opt(cfg)

	assert.Equal(t, float64(10), cfg.timeout.Seconds())
}

func TestWithDebug(t *testing.T) {
	cfg := &Options{
		debug: false,
	}

	opt := WithDebug(true)

	opt(cfg)

	assert.True(t, cfg.debug)
}

func TestWithLogger(t *testing.T) {
	cfg := &Options{
		logger: nil,
	}

	opt := WithLogger(&Log{})

	opt(cfg)

	assert.NotNil(t, cfg.logger)
}
