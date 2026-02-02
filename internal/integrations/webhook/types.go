package webhook

// EventType tipo de evento para webhook (compativel com WuzAPI)
type EventType string

// Tipos de eventos suportados
const (
	// Messages and Communication
	EventMessage              EventType = "Message"
	EventUndecryptableMessage EventType = "UndecryptableMessage"
	EventReceipt              EventType = "Receipt"
	EventMediaRetry           EventType = "MediaRetry"
	EventReadReceipt          EventType = "ReadReceipt"

	// Groups and Contacts
	EventGroupInfo       EventType = "GroupInfo"
	EventJoinedGroup     EventType = "JoinedGroup"
	EventPicture         EventType = "Picture"
	EventBlocklistChange EventType = "BlocklistChange"
	EventBlocklist       EventType = "Blocklist"

	// Connection and Session
	EventConnected         EventType = "Connected"
	EventDisconnected      EventType = "Disconnected"
	EventConnectFailure    EventType = "ConnectFailure"
	EventKeepAliveRestored EventType = "KeepAliveRestored"
	EventKeepAliveTimeout  EventType = "KeepAliveTimeout"
	EventQRTimeout         EventType = "QRTimeout"
	EventLoggedOut         EventType = "LoggedOut"
	EventClientOutdated    EventType = "ClientOutdated"
	EventTemporaryBan      EventType = "TemporaryBan"
	EventStreamError       EventType = "StreamError"
	EventStreamReplaced    EventType = "StreamReplaced"
	EventPairSuccess       EventType = "PairSuccess"
	EventPairError         EventType = "PairError"
	EventQR                EventType = "QR"

	// Privacy and Settings
	EventPrivacySettings EventType = "PrivacySettings"
	EventPushNameSetting EventType = "PushNameSetting"
	EventUserAbout       EventType = "UserAbout"

	// Synchronization and State
	EventAppState             EventType = "AppState"
	EventAppStateSyncComplete EventType = "AppStateSyncComplete"
	EventHistorySync          EventType = "HistorySync"
	EventOfflineSyncCompleted EventType = "OfflineSyncCompleted"
	EventOfflineSyncPreview   EventType = "OfflineSyncPreview"

	// Calls
	EventCallOffer        EventType = "CallOffer"
	EventCallAccept       EventType = "CallAccept"
	EventCallTerminate    EventType = "CallTerminate"
	EventCallOfferNotice  EventType = "CallOfferNotice"
	EventCallRelayLatency EventType = "CallRelayLatency"

	// Presence and Activity
	EventPresence     EventType = "Presence"
	EventChatPresence EventType = "ChatPresence"

	// Identity
	EventIdentityChange EventType = "IdentityChange"

	// Newsletter (WhatsApp Channels)
	EventNewsletterJoin       EventType = "NewsletterJoin"
	EventNewsletterLeave      EventType = "NewsletterLeave"
	EventNewsletterMuteChange EventType = "NewsletterMuteChange"
	EventNewsletterLiveUpdate EventType = "NewsletterLiveUpdate"

	// Facebook/Meta Bridge
	EventFBMessage EventType = "FBMessage"

	// Special - receives all events
	EventAll EventType = "All"
)

// Event representa um evento a ser enviado via webhook
type Event struct {
	Type      EventType   `json:"type"`
	SessionID string      `json:"sessionId"`
	Event     interface{} `json:"event,omitempty"`
}

// Config configuracao do webhook para uma sessao
type Config struct {
	URL        string      `json:"url"`
	Events     []EventType `json:"events"`
	HMACKeySet bool        `json:"hmacKeySet"`
}

// SupportedEvents retorna lista de tipos de eventos suportados
func SupportedEvents() []EventType {
	return []EventType{
		EventMessage,
		EventUndecryptableMessage,
		EventReceipt,
		EventMediaRetry,
		EventReadReceipt,
		EventGroupInfo,
		EventJoinedGroup,
		EventPicture,
		EventBlocklistChange,
		EventBlocklist,
		EventConnected,
		EventDisconnected,
		EventConnectFailure,
		EventKeepAliveRestored,
		EventKeepAliveTimeout,
		EventQRTimeout,
		EventLoggedOut,
		EventClientOutdated,
		EventTemporaryBan,
		EventStreamError,
		EventStreamReplaced,
		EventPairSuccess,
		EventPairError,
		EventQR,
		EventPrivacySettings,
		EventPushNameSetting,
		EventUserAbout,
		EventAppState,
		EventAppStateSyncComplete,
		EventHistorySync,
		EventOfflineSyncCompleted,
		EventOfflineSyncPreview,
		EventCallOffer,
		EventCallAccept,
		EventCallTerminate,
		EventCallOfferNotice,
		EventCallRelayLatency,
		EventPresence,
		EventChatPresence,
		EventIdentityChange,
		EventNewsletterJoin,
		EventNewsletterLeave,
		EventNewsletterMuteChange,
		EventNewsletterLiveUpdate,
		EventFBMessage,
		EventAll,
	}
}

// SupportedEventStrings retorna lista de eventos como strings
func SupportedEventStrings() []string {
	events := SupportedEvents()
	result := make([]string, len(events))
	for i, e := range events {
		result[i] = string(e)
	}
	return result
}

// ParseEventTypes converte strings para EventType
func ParseEventTypes(events []string) []EventType {
	validEvents := make(map[EventType]bool)
	for _, e := range SupportedEvents() {
		validEvents[e] = true
	}

	result := make([]EventType, 0, len(events))
	for _, e := range events {
		eventType := EventType(e)
		if validEvents[eventType] {
			result = append(result, eventType)
		}
	}
	return result
}

// EventTypesToStrings converte slice de EventType para strings
func EventTypesToStrings(events []EventType) []string {
	result := make([]string, len(events))
	for i, e := range events {
		result[i] = string(e)
	}
	return result
}
