package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
)

const CSRFSessionKey = "csrf_token"

type SessionManager struct {
	mgr *scs.SessionManager
}

var Session *SessionManager

func InitSessionManager() *scs.SessionManager {
	sessMgr := scs.New()
	sessMgr.Lifetime = 24 * time.Hour
	sessMgr.IdleTimeout = 2 * time.Hour
	sessMgr.Cookie.HttpOnly = true
	sessMgr.Cookie.SameSite = http.SameSiteLaxMode
	sessMgr.Cookie.Secure = strings.EqualFold(GetEnv("SESSION_COOKIE_SECURE", "false"), "true")
	Session = &SessionManager{mgr: sessMgr}
	return sessMgr
}

func (s *SessionManager) Get(r *http.Request, key string) interface{} {
	return s.mgr.Get(r.Context(), key)
}

func (s *SessionManager) Set(r *http.Request, key string, value interface{}) {
	s.mgr.Put(r.Context(), key, value)
}

func (s *SessionManager) RenewToken(r *http.Request) error {
	return s.mgr.RenewToken(r.Context())
}

func (s *SessionManager) Clear(r *http.Request) error {
	return s.mgr.Clear(r.Context())
}

func (s *SessionManager) CSRFToken(r *http.Request) (string, error) {
	value := s.Get(r, CSRFSessionKey)
	if token, ok := value.(string); ok && token != "" {
		return token, nil
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := base64.RawURLEncoding.EncodeToString(bytes)
	s.Set(r, CSRFSessionKey, token)
	return token, nil
}

func (s *SessionManager) VerifyCSRFToken(r *http.Request) bool {
	expected := s.Get(r, CSRFSessionKey)
	token, ok := expected.(string)
	if !ok || token == "" {
		return false
	}

	got := r.Header.Get("X-CSRF-Token")
	return subtle.ConstantTimeCompare([]byte(got), []byte(token)) == 1
}
