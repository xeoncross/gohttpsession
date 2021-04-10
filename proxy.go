package session

import (
	"github.com/xeoncross/gohttpsession/sessiontoken"
	"net/http"
)

// ID of a session
// type ID []byte
// We don't use this because we store the ID as a []byte slice
// (regardless of storage system: redis/mysql/mongo/dynamo/etc..)
// so using this alias would just result in more type conversions
// without providing type safety since we don't know the length
// that you will be using

// yeah, you thought I didn't see you there didn't you?

// Proxy instance for reading and writing sessions to HTTP clients
type Proxy struct {
	BaseCookie http.Cookie
	IDLength   int
}

// Load session ID (if exists)
func (p *Proxy) Load(r *http.Request) []byte {

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
func (p *Proxy) Start(w http.ResponseWriter) []byte {
	cookie := p.BaseCookie
	id := sessiontoken.New(p.IDLength)
	cookie.Value = sessiontoken.Encode(id)
	http.SetCookie(w, &cookie)
	return id
}
