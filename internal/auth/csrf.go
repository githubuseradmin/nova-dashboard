package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
)

const (
	// CSRFCookieName holds the CSRF token. It is readable by JS (not HttpOnly)
	// so the SPA can echo it back in a request header (double-submit pattern).
	CSRFCookieName = "nova_csrf"
	// CSRFHeaderName is the header the SPA sends the token back in.
	CSRFHeaderName = "X-CSRF-Token"
)

// NewCSRFToken returns a random, URL-safe CSRF token.
func NewCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// SetCSRFCookie writes the CSRF cookie so the SPA can read and echo it.
func SetCSRFCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false, // the SPA must read it to send X-CSRF-Token
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// ValidCSRF reports whether the request's CSRF header matches its CSRF cookie.
func ValidCSRF(r *http.Request) bool {
	cookie, err := r.Cookie(CSRFCookieName)
	if err != nil || cookie.Value == "" {
		return false
	}
	header := r.Header.Get(CSRFHeaderName)
	if header == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(header)) == 1
}
