package dto

// SendTextRequest request para enviar mensagem de texto
type SendTextRequest struct {
	Phone string `json:"Phone" example:"5511999999999"`
	Body  string `json:"Body" example:"Hello World!"`
}

// SendImageRequest request para enviar imagem
type SendImageRequest struct {
	Phone    string `json:"Phone" example:"5511999999999"`
	Image    string `json:"Image" example:"base64..."`
	Caption  string `json:"Caption,omitempty" example:"Image caption"`
	MimeType string `json:"Mimetype,omitempty" example:"image/jpeg"`
}

// SendVideoRequest request para enviar video
type SendVideoRequest struct {
	Phone    string `json:"Phone" example:"5511999999999"`
	Video    string `json:"Video" example:"base64..."`
	Caption  string `json:"Caption,omitempty" example:"Video caption"`
	MimeType string `json:"Mimetype,omitempty" example:"video/mp4"`
}

// SendDocumentRequest request para enviar documento
type SendDocumentRequest struct {
	Phone    string `json:"Phone" example:"5511999999999"`
	Document string `json:"Document" example:"base64..."`
	FileName string `json:"FileName" example:"document.pdf"`
	Caption  string `json:"Caption,omitempty" example:"Document caption"`
	MimeType string `json:"Mimetype,omitempty" example:"application/pdf"`
}

// SendAudioRequest request para enviar audio
type SendAudioRequest struct {
	Phone    string `json:"Phone" example:"5511999999999"`
	Audio    string `json:"Audio" example:"base64..."`
	MimeType string `json:"Mimetype,omitempty" example:"audio/ogg; codecs=opus"`
}

// SendStickerRequest request para enviar sticker
type SendStickerRequest struct {
	Phone    string `json:"Phone" example:"5511999999999"`
	Sticker  string `json:"Sticker" example:"base64..."`
	MimeType string `json:"Mimetype,omitempty" example:"image/webp"`
}

// SendLocationRequest request para enviar localizacao
type SendLocationRequest struct {
	Phone     string  `json:"Phone" example:"5511999999999"`
	Latitude  float64 `json:"Latitude" example:"-23.5505"`
	Longitude float64 `json:"Longitude" example:"-46.6333"`
	Name      string  `json:"Name,omitempty" example:"Sao Paulo"`
	Address   string  `json:"Address,omitempty" example:"Av. Paulista, 1000"`
}

// SendContactRequest request para enviar contato (vCard)
type SendContactRequest struct {
	Phone string `json:"Phone" example:"5511999999999"`
	Name  string `json:"Name" example:"John Doe"`
	Vcard string `json:"Vcard" example:"BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+5511999999999\nEND:VCARD"`
}

// SendPollRequest request para enviar enquete
type SendPollRequest struct {
	Phone       string   `json:"Phone" example:"5511999999999"`
	Question    string   `json:"Question" example:"What is your favorite color?"`
	Options     []string `json:"Options" example:"Red,Blue,Green"`
	MultiSelect bool     `json:"MultiSelect,omitempty" example:"false"`
}

// SendReactionRequest request para enviar reacao a mensagem
type SendReactionRequest struct {
	Phone     string `json:"Phone" example:"5511999999999"`
	MessageId string `json:"Id" example:"ABCD1234567890"`
	Emoji     string `json:"Body" example:"👍"`
}

// EditMessageRequest request para editar mensagem
type EditMessageRequest struct {
	Phone string `json:"Phone" example:"5511999999999"`
	Body  string `json:"Body" example:"Edited message text"`
}

// RevokeMessageRequest request para revogar/deletar mensagem
type RevokeMessageRequest struct {
	Phone string `json:"Phone" example:"5511999999999"`
}

// MessageResponse resposta com ID da mensagem enviada
type MessageResponse struct {
	MessageId string `json:"Id" example:"ABCD1234567890"`
	Timestamp int64  `json:"Timestamp,omitempty" example:"1704067200"`
}
