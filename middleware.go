package session

import (
	"context"
	"net/http"
)

type contextKey string

const ContextKey contextKey = "user"

type UserLookup = func([]byte) (interface{}, error)

func SetUserContext(next http.Handler, proxy *CookieProxy, load UserLookup) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionID := proxy.Load(r)
		if sessionID != nil {

			user, err := load(sessionID)
			if err == nil {
				ctx := context.WithValue(r.Context(), ContextKey, user)

				// Access context values in handlers like this
				// user, _ := r.Context().Value(usersession.ContextKey).(*MyUserType)

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
