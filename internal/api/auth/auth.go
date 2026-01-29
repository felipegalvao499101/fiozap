package auth

import (
	"context"
	"net/http"

	"fiozap/internal/api/dto"
	"fiozap/internal/domain"

	"github.com/go-chi/chi/v5"
)

type ctxKey string

const CtxKeyIsGlobal ctxKey = "isGlobal"

type Auth struct {
	globalToken string
	provider    domain.Provider
}

func NewAuth(globalToken string, provider domain.Provider) *Auth {
	return &Auth{
		globalToken: globalToken,
		provider:    provider,
	}
}

// Global verifica apenas o token global (para rotas admin como criar sessao)
func (a *Auth) Global(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			dto.Error(w, http.StatusUnauthorized, "missing token")
			return
		}

		if a.globalToken == "" || token != a.globalToken {
			dto.Error(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyIsGlobal, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Session verifica token global OU token da sessao especifica
func (a *Auth) Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			dto.Error(w, http.StatusUnauthorized, "missing token")
			return
		}

		// Token global tem acesso a tudo
		if a.globalToken != "" && token == a.globalToken {
			ctx := context.WithValue(r.Context(), CtxKeyIsGlobal, true)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Verifica token da sessao especifica
		name := chi.URLParam(r, "name")
		if name == "" {
			dto.Error(w, http.StatusBadRequest, "session name required")
			return
		}

		session, err := a.provider.GetSession(name)
		if err != nil {
			dto.Error(w, http.StatusNotFound, "session not found")
			return
		}

		if session.GetToken() != token {
			dto.Error(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyIsGlobal, false)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
