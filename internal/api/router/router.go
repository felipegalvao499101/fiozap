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
	groupHandler := handlers.NewGroupHandler(manager)
	chatHandler := handlers.NewChatHandler(manager)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/sessions", func(r chi.Router) {
		r.With(authMiddleware.Global).Post("/", sessionHandler.Create)
		r.With(authMiddleware.Global).Get("/", sessionHandler.List)

		r.Route("/{name}", func(r chi.Router) {
			r.Use(authMiddleware.Session)

			// Session
			r.Get("/", sessionHandler.Get)
			r.Post("/connect", sessionHandler.Connect)
			r.Get("/qr", sessionHandler.GetQR)
			r.Post("/disconnect", sessionHandler.Disconnect)
			r.Post("/logout", sessionHandler.Logout)
			r.Delete("/", sessionHandler.Delete)

			// Messages
			r.Route("/messages", func(r chi.Router) {
				r.Post("/text", messageHandler.SendText)
				r.Post("/image", messageHandler.SendImage)
				r.Post("/video", messageHandler.SendVideo)
				r.Post("/audio", messageHandler.SendAudio)
				r.Post("/document", messageHandler.SendDocument)
				r.Post("/sticker", messageHandler.SendSticker)
				r.Post("/location", messageHandler.SendLocation)
				r.Post("/contact", messageHandler.SendContact)
				r.Post("/poll", messageHandler.SendPoll)
				r.Post("/reaction", messageHandler.React)
				r.Put("/{messageId}", messageHandler.Edit)
				r.Delete("/{messageId}", messageHandler.Revoke)
			})

			// Users
			r.Post("/users/check", messageHandler.CheckPhone)
			r.Get("/users/{phone}", messageHandler.GetUserInfo)
			r.Get("/users/{phone}/avatar", messageHandler.GetUserAvatar)

			// Groups
			r.Route("/groups", func(r chi.Router) {
				r.Post("/", groupHandler.Create)
				r.Get("/", groupHandler.List)
				r.Post("/join", groupHandler.Join)
				r.Get("/invite/{code}", groupHandler.GetInviteInfo)

				r.Route("/{groupJid}", func(r chi.Router) {
					r.Get("/", groupHandler.Get)
					r.Put("/name", groupHandler.SetName)
					r.Put("/topic", groupHandler.SetTopic)
					r.Put("/photo", groupHandler.SetPhoto)
					r.Post("/leave", groupHandler.Leave)
					r.Get("/invite", groupHandler.GetInviteLink)
					r.Post("/invite/revoke", groupHandler.RevokeInviteLink)

					// Participants
					r.Route("/participants", func(r chi.Router) {
						r.Post("/", groupHandler.AddParticipants)
						r.Delete("/", groupHandler.RemoveParticipants)
						r.Post("/promote", groupHandler.PromoteParticipants)
						r.Post("/demote", groupHandler.DemoteParticipants)
					})

					// Settings
					r.Put("/settings/announce", groupHandler.SetAnnounce)
					r.Put("/settings/locked", groupHandler.SetLocked)
					r.Put("/settings/approval", groupHandler.SetApproval)
				})
			})

			// Chat
			r.Post("/chat/markread", chatHandler.MarkRead)
			r.Post("/chat/presence", chatHandler.Presence)
			r.Route("/chat/{chatJid}", func(r chi.Router) {
				r.Put("/disappearing", chatHandler.SetDisappearing)
			})

			// Presence (global)
			r.Post("/presence", chatHandler.SendPresence)
			r.Post("/presence/subscribe", chatHandler.SubscribePresence)

			// Blocklist
			r.Get("/blocklist", chatHandler.GetBlocklist)
			r.Post("/blocklist/block", chatHandler.BlockContact)
			r.Post("/blocklist/unblock", chatHandler.UnblockContact)
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
