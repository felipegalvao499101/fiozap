package dto

type SendTextRequest struct {
	To   string `json:"to" validate:"required"`
	Text string `json:"text" validate:"required"`
}

type SendImageRequest struct {
	To       string `json:"to" validate:"required"`
	Image    string `json:"image" validate:"required"`
	Caption  string `json:"caption,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

type SendDocumentRequest struct {
	To       string `json:"to" validate:"required"`
	Document string `json:"document" validate:"required"`
	Filename string `json:"filename" validate:"required"`
	MimeType string `json:"mime_type,omitempty"`
}

type SendAudioRequest struct {
	To       string `json:"to" validate:"required"`
	Audio    string `json:"audio" validate:"required"`
	MimeType string `json:"mime_type,omitempty"`
}

type SendLocationRequest struct {
	To        string  `json:"to" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

type MessageResponse struct {
	MessageID string `json:"message_id"`
}

type CheckPhoneRequest struct {
	Phones []string `json:"phones" validate:"required"`
}

type CheckPhoneResponse struct {
	Phone       string `json:"phone"`
	IsOnWhatsApp bool   `json:"is_on_whatsapp"`
	JID         string `json:"jid,omitempty"`
}

type UserInfoResponse struct {
	JID        string `json:"jid"`
	Name       string `json:"name,omitempty"`
	Status     string `json:"status,omitempty"`
	PictureURL string `json:"picture_url,omitempty"`
}
