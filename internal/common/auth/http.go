package auth

import (
	"net/http"
	"strings"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/common/server/httperr"
)

func BearerToken(r *http.Request) (string, error) {
	scheme, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if !ok || !strings.EqualFold(scheme, "bearer") || strings.TrimSpace(token) == "" {
		return "", commonerrors.NewAuthorizationError("missing bearer token", "unauthenticated")
	}
	return strings.TrimSpace(token), nil
}

// Middleware authenticates a request from the JWT signature and claims.
// It does not check that the session still exists. The user service does that
// on its own routes, so logout revokes access there immediately. Other services
// trust role and permissions from the token until it expires.
func Middleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := BearerToken(r)
			if err != nil {
				httperr.RespondWithSlugError(err, w, r)
				return
			}

			user, err := Parse(secret, token)
			if err != nil {
				httperr.RespondWithSlugError(err, w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
		})
	}
}
