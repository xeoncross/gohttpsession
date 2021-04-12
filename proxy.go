package session

import (
	"net/http"

	"github.com/xeoncross/gohttpsession/sessiontoken"
)

// CookieProxy for reading and writing session cookies to HTTP clients
type CookieProxy struct {
	BaseCookie http.Cookie
	IDLength   int
}

// Load session ID (if exists)
func (p *CookieProxy) Load(r *http.Request) []byte {

	sessionCookie, err := r.Cookie(p.BaseCookie.Name)
	if err != nil {
		return nil
	}

	if sessionCookie.Value != "" {
		sessionID := sessiontoken.Decode(sessionCookie.Value)
		if len(sessionID) == p.IDLength {
			return sessionID
		}
	}

	return nil
}

// Start a new session by sending a cookie with the new session ID to the client
func (p *CookieProxy) Start(w http.ResponseWriter) []byte {
	cookie := p.BaseCookie
	id := sessiontoken.New(p.IDLength)
	cookie.Value = sessiontoken.Encode(id)
	http.SetCookie(w, &cookie)
	return id
}
