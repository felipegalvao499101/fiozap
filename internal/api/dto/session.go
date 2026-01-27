package dto

type CreateSessionRequest struct {
	Name string `json:"Name"`
}

type SessionResponse struct {
	Name      string `json:"Name"`
	Token     string `json:"Token,omitempty"`
	JID       string `json:"JID,omitempty"`
	Phone     string `json:"Phone,omitempty"`
	PushName  string `json:"PushName,omitempty"`
	Connected bool   `json:"Connected"`
}

type QRResponse struct {
	QRCode string `json:"QRCode"`
}
