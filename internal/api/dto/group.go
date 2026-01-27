package dto

type CreateGroupRequest struct {
	Name         string   `json:"Name"`
	Participants []string `json:"Participants"`
}

type GroupResponse struct {
	JID          string                `json:"JID"`
	Name         string                `json:"Name"`
	Topic        string                `json:"Topic,omitempty"`
	GroupCreated string                `json:"GroupCreated,omitempty"`
	OwnerJID     string                `json:"OwnerJID,omitempty"`
	Participants []ParticipantResponse `json:"Participants,omitempty"`
	IsAnnounce   bool                  `json:"IsAnnounce,omitempty"`
	IsLocked     bool                  `json:"IsLocked,omitempty"`
}

type ParticipantResponse struct {
	JID          string `json:"JID"`
	IsAdmin      bool   `json:"IsAdmin"`
	IsSuperAdmin bool   `json:"IsSuperAdmin"`
}

type ParticipantsRequest struct {
	Phone []string `json:"Phone"`
}

type GroupNameRequest struct {
	Name string `json:"Name"`
}

type GroupDescriptionRequest struct {
	Topic string `json:"Topic"`
}

type GroupPhotoRequest struct {
	Image string `json:"Image"`
}

type JoinGroupRequest struct {
	Code string `json:"Code"`
}

type GroupSettingRequest struct {
	Value bool `json:"Value"`
}

type InviteLinkResponse struct {
	InviteLink string `json:"InviteLink"`
}
