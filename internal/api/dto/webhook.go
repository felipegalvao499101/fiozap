package dto

// SetWebhookRequest request para configurar webhook
type SetWebhookRequest struct {
	WebhookURL string   `json:"WebhookURL" example:"https://example.com/webhook"`
	Events     []string `json:"Events,omitempty" example:"Message,ReadReceipt,Connected"`
}

// SetHMACRequest request para configurar HMAC
type SetHMACRequest struct {
	HMACKey string `json:"HmacKey" example:"your_hmac_key_minimum_32_characters_long"`
}

// WebhookConfigResponse resposta com configuracao do webhook
type WebhookConfigResponse struct {
	Webhook    string   `json:"Webhook"`
	Events     []string `json:"Events"`
	HMACKeySet bool     `json:"HmacKeySet"`
}

// WebhookActionResponse resposta de acao do webhook
type WebhookActionResponse struct {
	Details string `json:"Details"`
}

// SupportedEventsResponse lista de eventos suportados
type SupportedEventsResponse struct {
	Events []string `json:"Events"`
}
