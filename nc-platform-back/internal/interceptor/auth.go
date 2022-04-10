package interceptor

import (
	"context"
	"github.com/go-chi/jwtauth"
	"net/http"
)

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", claims["userId"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
