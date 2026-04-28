package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"khairul169/garage-webui/utils"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct{}

const (
	maxFailedLoginAttempts = 5
	loginAttemptWindow     = 15 * time.Minute
	loginLockoutDuration   = 15 * time.Minute
)

type loginAttempt struct {
	failures     int
	firstFailure time.Time
	lockedUntil  time.Time
}

type loginRateLimiter struct {
	mu       sync.Mutex
	attempts map[string]loginAttempt
}

var authLimiter = &loginRateLimiter{attempts: map[string]loginAttempt{}}

func (c *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ResponseError(w, err)
		return
	}

	limiterKey := loginLimiterKey(r, body.Username)
	if retryAfter, locked := authLimiter.Check(limiterKey); locked {
		w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
		utils.ResponseErrorStatus(w, fmt.Errorf("too many failed login attempts; try again in %s", retryAfter.Round(time.Second)), http.StatusTooManyRequests)
		return
	}

	userPass := strings.SplitN(utils.GetEnv("AUTH_USER_PASS", ""), ":", 2)
	if len(userPass) < 2 {
		utils.ResponseErrorStatus(w, errors.New("AUTH_USER_PASS not set"), 500)
		return
	}

	if strings.TrimSpace(body.Username) != userPass[0] || bcrypt.CompareHashAndPassword([]byte(userPass[1]), []byte(body.Password)) != nil {
		if retryAfter, locked := authLimiter.RegisterFailure(limiterKey); locked {
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
			utils.ResponseErrorStatus(w, fmt.Errorf("too many failed login attempts; try again in %s", retryAfter.Round(time.Second)), http.StatusTooManyRequests)
			return
		}

		utils.ResponseErrorStatus(w, errors.New("invalid username or password"), 401)
		return
	}

	authLimiter.Reset(limiterKey)
	if err := utils.Session.RenewToken(r); err != nil {
		utils.ResponseError(w, err)
		return
	}

	utils.Session.Set(r, "authenticated", true)
	csrfToken, err := utils.Session.CSRFToken(r)
	if err != nil {
		utils.ResponseError(w, err)
		return
	}

	utils.ResponseSuccess(w, map[string]interface{}{
		"authenticated": true,
		"csrfToken":     csrfToken,
	})
}

func (c *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	utils.Session.Clear(r)
	utils.ResponseSuccess(w, true)
}

func (c *Auth) GetStatus(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := true
	authSession := utils.Session.Get(r, "authenticated")
	enabled := false

	if utils.GetEnv("AUTH_USER_PASS", "") != "" {
		enabled = true
	}

	if authSession != nil && authSession.(bool) {
		isAuthenticated = true
	}

	csrfToken, err := utils.Session.CSRFToken(r)
	if err != nil {
		utils.ResponseError(w, err)
		return
	}

	utils.ResponseSuccess(w, map[string]interface{}{
		"enabled":       enabled,
		"authenticated": isAuthenticated,
		"csrfToken":     csrfToken,
	})
}

func (l *loginRateLimiter) Check(key string) (time.Duration, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.cleanupLocked(now)

	attempt, ok := l.attempts[key]
	if !ok {
		return 0, false
	}

	if !attempt.lockedUntil.IsZero() && attempt.lockedUntil.After(now) {
		return time.Until(attempt.lockedUntil), true
	}

	if attempt.firstFailure.Add(loginAttemptWindow).Before(now) {
		delete(l.attempts, key)
	}

	return 0, false
}

func (l *loginRateLimiter) RegisterFailure(key string) (time.Duration, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.cleanupLocked(now)

	attempt := l.attempts[key]
	if attempt.firstFailure.IsZero() || attempt.firstFailure.Add(loginAttemptWindow).Before(now) {
		attempt = loginAttempt{firstFailure: now}
	}

	attempt.failures++
	if attempt.failures >= maxFailedLoginAttempts {
		attempt.lockedUntil = now.Add(loginLockoutDuration)
	}

	l.attempts[key] = attempt

	if !attempt.lockedUntil.IsZero() && attempt.lockedUntil.After(now) {
		return time.Until(attempt.lockedUntil), true
	}

	return 0, false
}

func (l *loginRateLimiter) cleanupLocked(now time.Time) {
	for key, attempt := range l.attempts {
		lockExpired := !attempt.lockedUntil.IsZero() && attempt.lockedUntil.Before(now)
		windowExpired := attempt.firstFailure.Add(loginAttemptWindow).Before(now)
		if lockExpired || windowExpired {
			delete(l.attempts, key)
		}
	}
}

func (l *loginRateLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, key)
}

func loginLimiterKey(r *http.Request, username string) string {
	return clientIP(r) + ":" + strings.ToLower(strings.TrimSpace(username))
}

func clientIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
