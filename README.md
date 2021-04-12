## Go HTTP Session

Ultra lightweight session handling functions for HTTP servers. Designed 
for server-side storage of session data with only a random token stored 
on a client cookie.

This library does not provide or assume a storage backend. It consists 
of a HTTP request cookie helper, a token generation, encoding, and decoding 
library, and storage agnostic middleware helper.

### Usage

If using HTTP cookies, create a HTTP cookie proxy for use by your handlers.

```go
proxy := &session.CookieProxy{
    BaseCookie: http.Cookie{
        Name:      "session",
        HttpOnly:  true,
        // Secure: true, // make sure to enable this in production!
    },
    IDLength: 32,
}
```

## Middleware

At this point you can use the middleware to automatically load the user object from the database using the session id (if present).

```go
mux := http.NewServeMux()
mux.Handle("/", adminHandler)

loaduser := func(sessionID []byte) (interface{}, error) {
	// load from database
	return &User{...}, nil
}

// Wrap context-setting middleware
handler := session.SetUserContext(mux, proxy, loaduser)

server := &http.Server{
	Addr:         ":80",
	Handler:      handler,
}

if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
    log.Fatal(err)
}

```

With the user record loaded, you can access it in your handlers

```go
handler := func(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value(session.ContextKey).(*User)
    if ok {
        // User record loaded, user logged in
    }
}
```


## Manual usage

You can also manually load or create the session cookies

```go
func privatepage(w http.ResponseWriter, r *http.Request) {
    sessionID := proxy.Load(r)
    if sessionID == nil {
    	// no session yet
    }
    ...
}

func login(w http.ResponseWriter, r *http.Request) {

    // Handle Login Here
	
    // Save or start session by creating session id and sending it to 
    // client in a cookie
    sessionID := proxy.Start(w)
    
    // Then store the user session data on the backend
    // i.e. store.SaveSession(sessionID, userDataHere)
}
```


## Tokens

You can also use the token library directly.

```go
// Create a new token
sessionID := sessiontoken.New(32)

// Base64 encode a token for use in URLs, Cookies, JSON, etc..
encodedToken := sessiontoken.Encode(sessionID)

// Decode a base64 token from a HTML form request (or JSON)
sessionID := sessiontoken.Decode(r.Form.Get("token"))
if sessionID != nil { 
	// valid
}
```

### Warning

It is highly recommended to transmit session tokens in a `HttpOnly` cookie 
with `Secure` set to true (over HTTPS) to help prevent theft of the cookie
by unauthorized client-side code (XSS).

To reduce the risk of Cross-Site Request Forgery (CSRF) it is also 
recommended that you use a library like https://github.com/justinas/nosurf
in addition to session handling.