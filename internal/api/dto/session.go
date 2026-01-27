package dto

type CreateSessionRequest struct {
	Name string `json:"name"`
}

type SessionResponse struct {
	Name      string `json:"name"`
	Token     string `json:"token,omitempty"`
	JID       string `json:"jid,omitempty"`
	Phone     string `json:"phone,omitempty"`
	PushName  string `json:"pushName,omitempty"`
	Connected bool   `json:"connected"`
}

type QRResponse struct {
	Code string `json:"code"`
}
