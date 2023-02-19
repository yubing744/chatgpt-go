package httpx

import (
	"bytes"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/publicsuffix"
)

// HttpSession 封装了 http.Client，实现了会话保持和 headers 的传递
type HttpSession struct {
	client *http.Client
}

// NewHttpSessionClient 返回一个新的 HttpSessionClient 实例
func NewHttpSession() (*HttpSession, error) {
	opts := &cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	cookieJar, err := cookiejar.New(opts)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Timeout: time.Second * 10,
		Jar:     cookieJar,
	}
	return &HttpSession{
		client: httpClient,
	}, nil
}

// Get 发送 GET 请求
func (httpx *HttpSession) Get(url string, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value[0])
	}

	resp, err := httpx.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Post 发送 POST 请求
func (httpx *HttpSession) Post(url string, headers http.Header, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value[0])
	}

	resp, err := httpx.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Cookies returns the value of a cookie
func (httpx *HttpSession) Cookies(host string) Coookies {
	domain := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/",
	}

	rawCookies := httpx.client.Jar.Cookies(domain)
	if rawCookies != nil {
		return Coookies(rawCookies)
	}

	return nil
}
