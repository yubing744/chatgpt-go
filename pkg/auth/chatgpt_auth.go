package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yubing744/chatgpt-go/pkg/httpx"
	"github.com/yubing744/chatgpt-go/pkg/utils"
)

// Error represents the base error class
type Error struct {
	location   string
	statusCode int
	details    string
}

func (e *Error) Error() string {
	return e.details
}

// Authenticator represents the OpenAI Authentication Reverse Engineered
type Authenticator struct {
	sessionToken string
	emailAddress string
	password     string
	proxy        string
	session      *httpx.HttpSession
	accessToken  string
	userAgent    string
}

// NewAuthenticator creates a new instance of Authenticator
func NewAuthenticator(emailAddress, password, proxy string) *Authenticator {
	auth := &Authenticator{
		emailAddress: emailAddress,
		password:     password,
		proxy:        proxy,
		userAgent:    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	}

	session, err := httpx.NewHttpSession(time.Second * 60)
	if err != nil {
		log.Fatal("init http session fail")
	}

	auth.session = session

	return auth
}

// urlEncode encodes the string to URL format
func urlEncode(str string) string {
	return url.QueryEscape(str)
}

// begin starts the authentication process
func (a *Authenticator) Begin() error {
	url := "https://explorer.api.openai.com/api/auth/csrf"
	headers := http.Header{
		"Host":            {"explorer.api.openai.com"},
		"Accept":          {"*/*"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {a.userAgent},
		"Accept-Language": {"en-GB,en-US;q=0.9,en;q=0.8"},
		"Referer":         {"https://explorer.api.openai.com/auth/login"},
		"Accept-Encoding": {"gzip, deflate, br"},
	}

	resp, err := a.session.Get(url, headers, true)
	if err != nil {
		return errors.Wrapf(err, "error in get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK &&
		strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var data struct {
			CsrfToken string `json:"csrfToken"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		err := a.partOne(data.CsrfToken)
		if err != nil {
			return err
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Error{
			location:   "Begin",
			statusCode: resp.StatusCode,
			details:    fmt.Sprintf("response error, detail: %s", body),
		}
	}

	return nil
}

func (a *Authenticator) partOne(token string) error {
	url := "https://explorer.api.openai.com/api/auth/signin/auth0?prompt=login"
	payload := fmt.Sprintf("callbackUrl=%s&csrfToken=%s&json=true", "%2F", token)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("User-Agent", a.userAgent)
	headers.Set("Host", "explorer.api.openai.com")
	headers.Set("Accept", "*/*")
	headers.Set("Accept-Language", "en-US,en;q=0.8")
	headers.Set("Origin", "https://explorer.api.openai.com")
	headers.Set("Referer", "https://explorer.api.openai.com/auth/login")
	headers.Set("Accept-Encoding", "gzip, deflate")

	resp, err := a.session.Post(url, headers, []byte(payload), true)
	if err != nil {
		return errors.Wrapf(err, "error in get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK &&
		strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var data struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}

		if data.URL == "https://explorer.api.openai.com/api/auth/error?error=OAuthSignin" || strings.Contains(data.URL, "error") {
			return &Error{
				location:   "partOne",
				statusCode: resp.StatusCode,
				details:    "You have been rate limited. Please try again later.",
			}
		}

		err := a.partTwo(data.URL)
		if err != nil {
			return err
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Error{
			location:   "partOne",
			statusCode: resp.StatusCode,
			details:    fmt.Sprintf("response error, detail: %s", body),
		}
	}

	return nil
}

func (a *Authenticator) partTwo(url string) error {
	headers := http.Header{
		"Host":            {"auth0.openai.com"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {a.userAgent},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Referer":         {"https://explorer.api.openai.com/"},
	}

	resp, err := a.session.Get(url, headers, true)
	if err != nil {
		return errors.Wrapf(err, "error in get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "error in read body in part three")
		}

		bodyString := string(body)
		state, ok := utils.RegexpExtra(bodyString, `state=([a-zA-Z0-9-_]*)`, 1)
		if !ok {
			return errors.New("not found state in respone body")
		}

		err = a.partThree(state)
		if err != nil {
			return err
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Error{
			location:   "partTwo",
			statusCode: resp.StatusCode,
			details:    fmt.Sprintf("response error, detail: %s", body),
		}
	}

	return nil
}

func (auth *Authenticator) partThree(state string) error {
	url := fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state)
	headers := http.Header{
		"Host":            []string{"auth0.openai.com"},
		"Accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Connection":      []string{"keep-alive"},
		"User-Agent":      []string{auth.userAgent},
		"Accept-Language": []string{"en-US,en;q=0.9"},
		"Referer":         []string{"https://explorer.api.openai.com/"},
	}

	resp, err := auth.session.Get(url, headers, true)
	if err != nil {
		return errors.Wrapf(err, "error in get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		err := auth.partFour(state)
		if err != nil {
			return err
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Error{
			location:   "partThree",
			statusCode: resp.StatusCode,
			details:    fmt.Sprintf("response error, detail: %s", body),
		}
	}

	return nil
}

func (a *Authenticator) partFour(state string) error {
	url := fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state)
	emailURLEncoded := urlEncode(a.emailAddress)

	headers := http.Header{}
	headers.Add("Host", "auth0.openai.com")
	headers.Add("Origin", "https://auth0.openai.com")
	headers.Add("Connection", "keep-alive")
	headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	headers.Add("User-Agent", a.userAgent)
	headers.Add("Referer", fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state))
	headers.Add("Accept-Language", "en-US,en;q=0.9")
	headers.Add("Content-Type", "application/x-www-form-urlencoded")

	payload := fmt.Sprintf("state=%s&username=%s&js-available=false&webauthn-available=true&is-brave=false&webauthn-platform-available=true&action=default", state, emailURLEncoded)

	resp, err := a.session.Post(url, headers, []byte(payload), true)
	if err != nil {
		return errors.Wrapf(err, "error in get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 302 || resp.StatusCode == 200 {
		err = a.partFive(state)
		if err != nil {
			return err
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Error{
			location:   "partFour",
			statusCode: resp.StatusCode,
			details:    fmt.Sprintf("response error, detail: %s", body),
		}
	}

	return nil
}

func (a *Authenticator) partFive(state string) error {
	url := fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state)
	headers := http.Header{
		"Host":            {"auth0.openai.com"},
		"Origin":          {"https://auth0.openai.com"},
		"Connection":      {"keep-alive"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"User-Agent":      {a.userAgent},
		"Referer":         {fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state)},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Content-Type":    {"application/x-www-form-urlencoded"},
	}

	emailURLEncoded := urlEncode(a.emailAddress)
	passwordURLEncoded := urlEncode(a.password)
	payload := fmt.Sprintf("state=%s&username=%s&password=%s&action=default", state, emailURLEncoded, passwordURLEncoded)

	resp, err := a.session.Post(url, headers, []byte(payload), false)
	if err != nil {
		return errors.Wrapf(err, "error in get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 302 || resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "error in read body in part five")
		}

		bodyString := string(body)
		newState, ok := utils.RegexpExtra(bodyString, `state=([a-zA-Z0-9-_]*)`, 1)
		if !ok {
			fmt.Print(bodyString)
			return errors.New("not found state in respone body of part five")
		}

		err = a.partSix(state, newState)
		if err != nil {
			return err
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Error{
			location:   "partFive",
			statusCode: resp.StatusCode,
			details:    fmt.Sprintf("response error, detail: %s", body),
		}
	}

	return nil
}

func (a *Authenticator) partSix(oldState, newState string) error {
	url := fmt.Sprintf("https://auth0.openai.com/authorize/resume?state=%s", newState)

	headers := http.Header{
		"Host":            {"auth0.openai.com"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {a.userAgent},
		"Accept-Language": {"en-GB,en-US;q=0.9,en;q=0.8"},
		"Referer":         {fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", oldState)},
	}

	resp, err := a.session.Get(url, headers, false)
	if err != nil {
		return errors.Wrapf(err, "error in get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		// Print redirect url
		redirectURL := resp.Header.Get("Location")
		if err = a.partSeven(redirectURL, url); err != nil {
			return err
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Error{
			location:   "partSix",
			statusCode: resp.StatusCode,
			details:    fmt.Sprintf("response error, detail: %s", body),
		}
	}

	return nil
}

func (a *Authenticator) partSeven(redirectURL string, previousURL string) error {
	url := redirectURL
	headers := http.Header{
		"Host":            {"explorer.api.openai.com"},
		"Accept":          {"application/json"},
		"Connection":      {"keep-alive"},
		"User-Agent":      {a.userAgent},
		"Accept-Language": {"en-GB,en-US;q=0.9,en;q=0.8"},
		"Referer":         {previousURL},
	}

	resp, err := a.session.Get(url, headers, false)
	if err != nil {
		return errors.Wrapf(err, "error in get %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		cookies := httpx.Coookies(resp.Cookies())
		sessionToken, ok := cookies.Get("__Secure-next-auth.session-token")

		if ok {
			a.sessionToken = sessionToken
			_, err = a.GetAccessToken()
			if err != nil {
				return err
			}
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return &Error{
			location:   "partSeven",
			statusCode: resp.StatusCode,
			details:    fmt.Sprintf("response error, detail: %s", body),
		}
	}

	return nil
}

func (a *Authenticator) GetSessionToken() string {
	return a.sessionToken
}

func (a *Authenticator) GetAccessToken() (string, error) {
	a.session.Cookies("openai.com").Set(
		"__Secure-next-auth.session-token",
		a.sessionToken,
	)

	resp, err := a.session.Get("https://explorer.api.openai.com/api/auth/session", nil, true)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 &&
		strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var data struct {
			AccessToken string `json:"accessToken"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return "", err
		}
		a.accessToken = data.AccessToken
		return a.accessToken, nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", &Error{
			location:   "GetAccessToken",
			statusCode: resp.StatusCode,
			details:    fmt.Sprintf("response error, detail: %s", body),
		}
	}
}
