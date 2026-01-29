package dto

// MarkReadRequest request para marcar mensagens como lidas
type MarkReadRequest struct {
	Id          []string `json:"Id" example:"ABCD1234567890,EFGH0987654321"`
	ChatPhone   string   `json:"ChatPhone,omitempty" example:"5511999999999"`
	SenderPhone string   `json:"SenderPhone,omitempty" example:"5511888888888"`
}

// ChatPresenceRequest request para enviar presenca no chat (digitando/gravando)
type ChatPresenceRequest struct {
	Phone string `json:"Phone" example:"5511999999999"`
	State string `json:"State" example:"composing" enums:"composing,paused"`
	Media string `json:"Media,omitempty" example:"" enums:",audio"`
}

// DisappearingRequest request para configurar mensagens temporarias
type DisappearingRequest struct {
	Duration string `json:"Duration" example:"24h" enums:"24h,7d,90d,off"`
}

// PresenceRequest request para enviar presenca global (online/offline)
type PresenceRequest struct {
	Type string `json:"Type" example:"available" enums:"available,unavailable"`
}

// PresenceSubscribeRequest request para se inscrever em presenca de contato
type PresenceSubscribeRequest struct {
	Phone string `json:"Phone" example:"5511999999999"`
}

// ChatActionResponse resposta generica de acoes de chat
type ChatActionResponse struct {
	Details string `json:"Details" example:"Messages marked as read"`
}
