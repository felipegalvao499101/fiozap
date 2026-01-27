package zap

import (
	"github.com/rs/zerolog"
	"go.mau.fi/whatsmeow/types/events"
)

func handleEvent(session *Session, log zerolog.Logger, evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		session.SetConnected(true)
		session.SetQRCode("")
		log.Info().Str("name", session.Name).Msg("Connected")

	case *events.Disconnected:
		session.SetConnected(false)
		log.Info().Str("name", session.Name).Msg("Disconnected")

	case *events.LoggedOut:
		session.SetConnected(false)
		log.Info().Str("name", session.Name).Str("reason", v.Reason.String()).Msg("Logged out")

	case *events.Message:
		log.Debug().
			Str("from", v.Info.Sender.String()).
			Str("chat", v.Info.Chat.String()).
			Msg("Message received")

	case *events.Receipt:
		log.Debug().
			Str("type", string(v.Type)).
			Msg("Receipt")

	case *events.HistorySync:
		log.Debug().Msg("History sync")

	case *events.PushName:
		log.Debug().Str("name", v.NewPushName).Msg("Push name update")
	}
}
