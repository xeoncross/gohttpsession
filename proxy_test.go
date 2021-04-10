package session_test

import (
	session "github.com/xeoncross/gohttpsession"
	"github.com/xeoncross/gohttpsession/sessiontoken"
	"net/http"
	"net/http/httptest"
	"testing"

)

func TestProxy(t *testing.T) {

	proxy := &session.Proxy{
		BaseCookie: http.Cookie{
			Name:     "session",
			HttpOnly: true,
			// Secure:   true, // make sure to enable this in production!
		},
		IDLength: 32,
	}

	rr := httptest.NewRecorder()
	id := proxy.Start(rr)

	// Parse cookie header back into cookie object
	rawCookie := rr.Result().Header.Get("set-cookie")
	header := http.Header{}
	header.Add("Cookie", rawCookie)
	request := http.Request{Header: header}
	cookie := request.Cookies()[0]

	if cookie.Value != sessiontoken.Encode(id) {
		t.Fatal("failed to set session cookie value")
	}

	if cookie.Name != "session" {
		t.Fatal("failed to set session cookie name")
	}

}