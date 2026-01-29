package dto

// CheckPhoneRequest request para verificar numeros no WhatsApp
type CheckPhoneRequest struct {
	Phone []string `json:"Phone" example:"5511999999999,5521888888888"`
}

// CheckPhoneResponse resposta da verificacao de numero
type CheckPhoneResponse struct {
	Query        string `json:"Query" example:"5511999999999"`
	IsInWhatsapp bool   `json:"IsInWhatsapp" example:"true"`
	JID          string `json:"JID,omitempty" example:"5511999999999@s.whatsapp.net"`
}

// UserInfoResponse informacoes do usuario/contato
type UserInfoResponse struct {
	JID       string   `json:"JID" example:"5511999999999@s.whatsapp.net"`
	Status    string   `json:"Status,omitempty" example:"Hey there! I am using WhatsApp."`
	PictureID string   `json:"PictureID,omitempty"`
	Devices   []string `json:"Devices,omitempty"`
}

// AvatarResponse resposta com foto de perfil
type AvatarResponse struct {
	URL string `json:"URL" example:"https://pps.whatsapp.net/..."`
	ID  string `json:"ID" example:"122345678901234567890"`
}

// ContactsCheckResponse resposta do check de contatos
type ContactsCheckResponse struct {
	Contacts []CheckPhoneResponse `json:"Contacts"`
}
