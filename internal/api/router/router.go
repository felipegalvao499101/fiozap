package router

import (
	"net/http"
	"time"

	"github.com/fiozap/fiozap/internal/api/handlers"
	"github.com/fiozap/fiozap/internal/api/middleware/auth"
	"github.com/fiozap/fiozap/internal/zap"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func New(manager *zap.Manager, logger zerolog.Logger, globalToken string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(requestLogger(logger))
	r.Use(middleware.Timeout(60 * time.Second))

	authMiddleware := auth.New(globalToken, manager)
	sessionHandler := handlers.NewSessionHandler(manager)
	messageHandler := handlers.NewMessageHandler(manager)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/sessions", func(r chi.Router) {
		// Rotas que requerem token global
		r.With(authMiddleware.Global).Post("/", sessionHandler.Create)
		r.With(authMiddleware.Global).Get("/", sessionHandler.List)

		// Rotas que aceitam token global ou token da sessao
		r.Route("/{name}", func(r chi.Router) {
			r.Use(authMiddleware.Session)

			r.Get("/", sessionHandler.Get)
			r.Post("/connect", sessionHandler.Connect)
			r.Get("/qr", sessionHandler.GetQR)
			r.Post("/disconnect", sessionHandler.Disconnect)
			r.Delete("/", sessionHandler.Delete)

			r.Post("/messages/text", messageHandler.SendText)
			r.Post("/messages/image", messageHandler.SendImage)
			r.Post("/messages/document", messageHandler.SendDocument)
			r.Post("/messages/audio", messageHandler.SendAudio)
			r.Post("/messages/location", messageHandler.SendLocation)

			r.Post("/users/check", messageHandler.CheckPhone)
			r.Get("/users/{phone}", messageHandler.GetUserInfo)
			r.Get("/users/{phone}/avatar", messageHandler.GetUserAvatar)
		})
	})

	return r
}

func requestLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", ww.Status()).
				Dur("duration", time.Since(start)).
				Msg("request")
		})
	}
}
