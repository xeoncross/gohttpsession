package session_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	session "github.com/xeoncross/gohttpsession"
	"github.com/xeoncross/gohttpsession/sessiontoken"
)

func TestUserContext(t *testing.T) {

	// Create the sesssion proxy
	proxy := &session.CookieProxy{
		BaseCookie: http.Cookie{
			Name:     "session",
			HttpOnly: true,
			// Secure:   true, // make sure to enable this in production!
		},
		IDLength: 24,
	}

	// Create a new session
	sessionID := sessiontoken.New(proxy.IDLength)

	// and the cookie we will set for this request for the middleware
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessiontoken.Encode(sessionID),
	}

	// Create the request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(cookie)

	rr := httptest.NewRecorder()

	// The handler that will load the user record
	handler := func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(session.ContextKey).(*User)
		if !ok {
			fmt.Fprintf(w, "unexpected: %+v\n", user)
		} else {
			fmt.Fprint(w, user.Name)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handler))

	// Wrap context-setting middleware
	h := session.SetUserContext(mux, proxy, loaduser)

	h.ServeHTTP(rr, req)

	if rr.Body.String() != "John" {
		t.Errorf("invalid response: %s", rr.Body.String())
	}

}

type User struct {
	Name string
}

// would load the user record from the database
func loaduser(sessionID []byte) (interface{}, error) {
	// return nil, errors.New("problem")
	return &User{"John"}, nil
}
