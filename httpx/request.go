package httpx

import (
	"io"
	"net/http"
	"strings"
	"time"
)

func RequestHeaders(request *http.Request) map[string]any {
	headers := make(map[string]any, len(request.Header))
	for key, value := range request.Header {
		headers[key] = strings.Join(value, ",")
	}
	return headers
}

func RequestHeader(request *http.Request, key string) Param {
	return NewParam(request.Header.Get(key))
}

func RequestCookies(request *http.Request) map[string]any {
	cookies := make(map[string]any)
	for _, cookie := range request.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}

func RequestCookie(request *http.Request, name string) Param {
	cookie, err := request.Cookie(name)
	if err != nil {
		return EmptyParam()
	}

	return NewParam(cookie.Value)
}

func NewCookie(key, value string, ttl ...time.Duration) *http.Cookie {
	const defaultExpire = time.Hour * 24 * 7

	expireAt := time.Now().Add(defaultExpire)
	if len(ttl) > 0 && ttl[0] > 0 {
		expireAt = time.Now().Add(ttl[0])
	}

	cookie := &http.Cookie{}
	cookie.Name = key
	cookie.Value = value
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Expires = expireAt
	return cookie
}

func RequestBody(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, nil
	}

	return io.ReadAll(request.Body)
}
