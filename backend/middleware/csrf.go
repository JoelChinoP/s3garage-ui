package middleware

import (
	"errors"
	"khairul169/garage-webui/utils"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isSafeMethod(r.Method) {
			next.ServeHTTP(w, r)
			return
		}

		if !isSameOrigin(r) {
			utils.ResponseErrorStatus(w, errors.New("invalid request origin"), http.StatusForbidden)
			return
		}

		if !utils.Session.VerifyCSRFToken(r) {
			utils.ResponseErrorStatus(w, errors.New("invalid csrf token"), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

func isSameOrigin(r *http.Request) bool {
	requestHost := normalizedHost(r.Host)
	if requestHost == "" {
		return false
	}

	if origin := r.Header.Get("Origin"); origin != "" {
		return urlHostMatches(origin, requestHost)
	}

	if referer := r.Header.Get("Referer"); referer != "" {
		return urlHostMatches(referer, requestHost)
	}

	return true
}

func urlHostMatches(rawURL string, requestHost string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return normalizedHost(parsed.Host) == requestHost
}

func normalizedHost(host string) string {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "" {
		return ""
	}

	if hostname, port, err := net.SplitHostPort(host); err == nil {
		if port == "80" || port == "443" {
			return hostname
		}
		return hostname + ":" + port
	}

	return strings.TrimSuffix(host, ".")
}
