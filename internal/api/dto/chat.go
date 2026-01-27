package dto

type MarkReadRequest struct {
	Id          []string `json:"Id"`
	ChatPhone   string   `json:"ChatPhone,omitempty"`
	SenderPhone string   `json:"SenderPhone,omitempty"`
}

type ChatPresenceRequest struct {
	Phone string `json:"Phone"`
	State string `json:"State"`
	Media string `json:"Media,omitempty"`
}

type DisappearingRequest struct {
	Duration string `json:"Duration"`
}

// Presence
type PresenceRequest struct {
	Type string `json:"Type"`
}

type PresenceSubscribeRequest struct {
	Phone string `json:"Phone"`
}

// Blocklist
type BlocklistResponse struct {
	JIDs []string `json:"JIDs"`
}

type BlockRequest struct {
	Phone string `json:"Phone"`
}
