package httpx

import "net/http"

type Coookies []*http.Cookie

func (c Coookies) Get(name string) (string, bool) {
	for _, item := range c {
		if item.Name == name {
			return item.Value, true
		}
	}

	return "", false
}

func (c Coookies) Set(name string, val string) bool {
	for _, item := range c {
		if item.Name == name {
			item.Value = val
			return true
		}
	}

	return false
}
