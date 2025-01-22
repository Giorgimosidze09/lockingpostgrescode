package middleware

import (
	"net/http"

	"lockingpostgrescode/auth"
)

func RoleMiddleware(requiredRole string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(auth.Claims)
		if !ok || (claims.Role != requiredRole && claims.Role != "super_admin") {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
