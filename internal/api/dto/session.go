package dto

// ProviderType tipo do provider
type ProviderType string

const (
	ProviderWhatsmeow ProviderType = "whatsmeow" // default - nao-oficial
	ProviderCloudAPI  ProviderType = "cloudapi"  // API oficial Meta
)

type CreateSessionRequest struct {
	Name     string       `json:"Name"`
	Provider ProviderType `json:"Provider,omitempty"` // "whatsmeow" (default) ou "cloudapi"
}

type SessionResponse struct {
	Name      string       `json:"Name"`
	Token     string       `json:"Token,omitempty"`
	JID       string       `json:"JID,omitempty"`
	Phone     string       `json:"Phone,omitempty"`
	PushName  string       `json:"PushName,omitempty"`
	Connected bool         `json:"Connected"`
	Provider  ProviderType `json:"Provider,omitempty"`
}

type QRResponse struct {
	QRCode string `json:"QRCode"`
}
