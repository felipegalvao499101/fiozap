package dto

type SendTextRequest struct {
	Phone string `json:"Phone"`
	Body  string `json:"Body"`
}

type SendImageRequest struct {
	Phone    string `json:"Phone"`
	Image    string `json:"Image"`
	Caption  string `json:"Caption,omitempty"`
	MimeType string `json:"Mimetype,omitempty"`
}

type SendVideoRequest struct {
	Phone    string `json:"Phone"`
	Video    string `json:"Video"`
	Caption  string `json:"Caption,omitempty"`
	MimeType string `json:"Mimetype,omitempty"`
}

type SendDocumentRequest struct {
	Phone    string `json:"Phone"`
	Document string `json:"Document"`
	FileName string `json:"FileName"`
	Caption  string `json:"Caption,omitempty"`
	MimeType string `json:"Mimetype,omitempty"`
}

type SendAudioRequest struct {
	Phone    string `json:"Phone"`
	Audio    string `json:"Audio"`
	MimeType string `json:"Mimetype,omitempty"`
}

type SendStickerRequest struct {
	Phone    string `json:"Phone"`
	Sticker  string `json:"Sticker"`
	MimeType string `json:"Mimetype,omitempty"`
}

type SendLocationRequest struct {
	Phone     string  `json:"Phone"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	Name      string  `json:"Name,omitempty"`
	Address   string  `json:"Address,omitempty"`
}

type SendContactRequest struct {
	Phone string `json:"Phone"`
	Name  string `json:"Name"`
	Vcard string `json:"Vcard"`
}

type SendPollRequest struct {
	Phone       string   `json:"Phone"`
	Question    string   `json:"Question"`
	Options     []string `json:"Options"`
	MultiSelect bool     `json:"MultiSelect,omitempty"`
}

type SendReactionRequest struct {
	Phone     string `json:"Phone"`
	MessageId string `json:"Id"`
	Emoji     string `json:"Body"`
}

type EditMessageRequest struct {
	Phone string `json:"Phone"`
	Body  string `json:"Body"`
}

type RevokeMessageRequest struct {
	Phone string `json:"Phone"`
}

type MessageResponse struct {
	MessageId string `json:"Id"`
	Timestamp int64  `json:"Timestamp,omitempty"`
}

type CheckPhoneRequest struct {
	Phone []string `json:"Phone"`
}

type CheckPhoneResponse struct {
	Query        string `json:"Query"`
	IsInWhatsapp bool   `json:"IsInWhatsapp"`
	JID          string `json:"JID,omitempty"`
}

type UserInfoResponse struct {
	JID        string   `json:"JID"`
	Status     string   `json:"Status,omitempty"`
	PictureID  string   `json:"PictureID,omitempty"`
	Devices    []string `json:"Devices,omitempty"`
}
