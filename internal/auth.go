package internal

import (
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware: extract username from JWT, put in context
func AuthMiddleware(userSvc *UserService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		username, err := userSvc.ParseToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "username", username)
		next(w, r.WithContext(ctx))
	}
}

// UsernameFromCtx extracts username from context
func UsernameFromCtx(ctx context.Context) string {
	v := ctx.Value("username")
	if v == nil {
		return ""
	}
	return v.(string)
}
