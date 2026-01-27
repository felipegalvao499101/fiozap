package auth

import (
	"context"
	"net/http"

	"github.com/fiozap/fiozap/internal/api/dto"
	"github.com/fiozap/fiozap/internal/zap"
	"github.com/go-chi/chi/v5"
)

type ctxKey string

const CtxKeyIsGlobal ctxKey = "isGlobal"

type Auth struct {
	globalToken string
	manager     *zap.Manager
}

func New(globalToken string, manager *zap.Manager) *Auth {
	return &Auth{
		globalToken: globalToken,
		manager:     manager,
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

		session, err := a.manager.GetSession(name)
		if err != nil {
			dto.Error(w, http.StatusNotFound, "session not found")
			return
		}

		if session.Token != token {
			dto.Error(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyIsGlobal, false)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
