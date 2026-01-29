package router

import (
	"net/http"
	"time"

	"fiozap/internal/api/auth"
	"fiozap/internal/api/handlers"
	"fiozap/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(provider domain.Provider, logger zerolog.Logger, globalToken string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(requestLogger(logger))
	r.Use(middleware.Timeout(60 * time.Second))

	authMiddleware := auth.NewAuth(globalToken, provider)
	sessionHandler := handlers.NewSessionHandler(provider)
	messageHandler := handlers.NewMessageHandler(provider)
	contactHandler := handlers.NewContactHandler(provider)
	groupHandler := handlers.NewGroupHandler(provider)
	chatHandler := handlers.NewChatHandler(provider)
	blocklistHandler := handlers.NewBlocklistHandler(provider)
	callHandler := handlers.NewCallHandler(provider)
	newsletterHandler := handlers.NewNewsletterHandler(provider)
	privacyHandler := handlers.NewPrivacyHandler(provider)
	profileHandler := handlers.NewProfileHandler(provider)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

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

			// Contacts
			r.Post("/contacts/check", contactHandler.CheckPhone)
			r.Get("/contacts/{phone}", contactHandler.GetInfo)
			r.Get("/contacts/{phone}/avatar", contactHandler.GetAvatar)
			r.Get("/contacts/{phone}/business", contactHandler.GetBusinessProfile)

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
			r.Get("/blocklist", blocklistHandler.GetBlocklist)
			r.Post("/blocklist/block", blocklistHandler.Block)
			r.Post("/blocklist/unblock", blocklistHandler.Unblock)

			// Newsletter (Channels)
			r.Route("/newsletters", func(r chi.Router) {
				r.Post("/", newsletterHandler.Create)
				r.Get("/", newsletterHandler.List)
				r.Route("/{newsletterJid}", func(r chi.Router) {
					r.Get("/", newsletterHandler.Get)
					r.Post("/follow", newsletterHandler.Follow)
					r.Post("/unfollow", newsletterHandler.Unfollow)
					r.Put("/mute", newsletterHandler.ToggleMute)
					r.Post("/reaction", newsletterHandler.SendReaction)
				})
			})

			// Privacy
			r.Route("/privacy", func(r chi.Router) {
				r.Get("/", privacyHandler.GetSettings)
				r.Put("/", privacyHandler.SetSetting)
				r.Get("/status", privacyHandler.GetStatusPrivacy)
			})

			// Profile
			r.Route("/profile", func(r chi.Router) {
				r.Get("/qrlink", profileHandler.GetContactQRLink)
				r.Post("/qrlink/resolve", profileHandler.ResolveContactQRLink)
				r.Put("/status", profileHandler.SetStatusMessage)
				r.Post("/business/resolve", profileHandler.ResolveBusinessMessageLink)
			})

			// Group Request Participants
			r.Get("/groups/{groupJid}/requests", groupHandler.GetRequestParticipants)
			r.Post("/groups/{groupJid}/requests/approve", groupHandler.ApproveRequestParticipants)
			r.Post("/groups/{groupJid}/requests/reject", groupHandler.RejectRequestParticipants)
			r.Put("/groups/{groupJid}/settings/memberadd", groupHandler.SetMemberAddMode)

			// Community
			r.Post("/community/link", groupHandler.LinkGroup)
			r.Post("/community/unlink", groupHandler.UnlinkGroup)
			r.Get("/community/{communityJid}/subgroups", groupHandler.GetSubGroups)
			r.Get("/community/{communityJid}/participants", groupHandler.GetLinkedParticipants)

			// Calls
			r.Post("/calls/reject", callHandler.RejectCall)
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
