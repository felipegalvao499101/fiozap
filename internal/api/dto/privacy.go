package dto

// Privacy Settings

type PrivacySettingsResponse struct {
	GroupAdd       string `json:"GroupAdd,omitempty"`       // all, contacts, contact_blacklist
	LastSeen       string `json:"LastSeen,omitempty"`       // all, contacts, contact_blacklist, none
	Status         string `json:"Status,omitempty"`         // all, contacts, contact_blacklist, none
	Profile        string `json:"Profile,omitempty"`        // all, contacts, contact_blacklist, none
	ReadReceipts   string `json:"ReadReceipts,omitempty"`   // all, none
	CallAdd        string `json:"CallAdd,omitempty"`        // all, known
	Online         string `json:"Online,omitempty"`         // all, match_last_seen
}

type SetPrivacyRequest struct {
	Name  string `json:"Name"`  // groupAdd, lastSeen, status, profile, readReceipts, callAdd, online
	Value string `json:"Value"` // all, contacts, contact_blacklist, none, known, match_last_seen
}

type StatusPrivacyResponse struct {
	Type      string   `json:"Type"`
	JIDList   []string `json:"JIDList,omitempty"`
	IsDefault bool     `json:"IsDefault"`
}
