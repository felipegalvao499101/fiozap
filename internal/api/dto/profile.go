package dto

// Profile

type SetStatusMessageRequest struct {
	Status string `json:"Status"`
}

type ContactQRLinkResponse struct {
	Link string `json:"Link"`
}

type ResolveContactQRRequest struct {
	Link string `json:"Link"`
}

type ContactQRTargetResponse struct {
	JID      string `json:"JID"`
	PushName string `json:"PushName,omitempty"`
}

type BusinessProfileResponse struct {
	JID         string   `json:"JID"`
	Description string   `json:"Description,omitempty"`
	Email       string   `json:"Email,omitempty"`
	Website     []string `json:"Website,omitempty"`
	Category    string   `json:"Category,omitempty"`
	Address     string   `json:"Address,omitempty"`
}

type ResolveBusinessLinkRequest struct {
	Link string `json:"Link"`
}

type BusinessLinkTargetResponse struct {
	JID     string `json:"JID"`
	Message string `json:"Message,omitempty"`
}
