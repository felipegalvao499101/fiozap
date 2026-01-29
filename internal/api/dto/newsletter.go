package dto

// Newsletter

type CreateNewsletterRequest struct {
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	Picture     string `json:"Picture,omitempty"` // base64
}

type NewsletterResponse struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description,omitempty"`
	Subscribers int    `json:"Subscribers,omitempty"`
	Verified    bool   `json:"Verified,omitempty"`
	Muted       bool   `json:"Muted,omitempty"`
	InviteLink  string `json:"InviteLink,omitempty"`
}

type NewsletterMessageRequest struct {
	Body string `json:"Body"`
}

type NewsletterReactionRequest struct {
	ServerID string `json:"ServerID"`
	Reaction string `json:"Reaction"` // emoji ou vazio para remover
}

type NewsletterMuteRequest struct {
	Mute bool `json:"Mute"`
}
