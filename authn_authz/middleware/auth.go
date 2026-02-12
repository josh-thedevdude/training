package middleware

import (
	"authn_authz/helpers"
	"context"
	"net/http"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the authorization header
		authHeader := r.Header.Get("authorization")

		if authHeader == "" {
			err := &helpers.AppError{
				Code:    401,
				Message: "unauthorized access",
				Tag:     "missing_authorization_header",
			}
			http.Error(w, err.Message, err.Code)
			return
		}

		// verify the authorization header
		claims, err := helpers.VerifyJwtToken(authHeader)
		if err != nil {
			if appErr, ok := err.(*helpers.AppError); ok {
				http.Error(w, appErr.Message, appErr.Code)
				return
			}
			http.Error(w, "unknown error", 500)
			return
		}

		// set the request context

		ctx := context.WithValue(r.Context(), helpers.UserIDKey, claims.UserId)
		ctx = context.WithValue(ctx, helpers.RoleKey, claims.Role)

		// attach context to the request
		r = r.WithContext(ctx)

		// call the next handler
		next.ServeHTTP(w, r)
	}
}

func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(helpers.RoleKey).(string)
			if !ok {
				http.Error(w, "unauthorized", 401)
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next(w, r)
					return
				}
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
}

// func middleware(originalHandler http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         fmt.Println("Running before handler")
//         w.Write([]byte("Hijacking Request "))
//         originalHandler.Serve(w, r)
//         fmt.Println("Running after handler")
//     })
// }
