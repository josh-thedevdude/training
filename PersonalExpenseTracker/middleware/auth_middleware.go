package middleware

import "net/http"

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc()
}
