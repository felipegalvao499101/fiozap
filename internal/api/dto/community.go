package dto

// Community (grupos ligados)

type LinkGroupRequest struct {
	ParentJID string `json:"ParentJID"`
	ChildJID  string `json:"ChildJID"`
}

type SubGroupResponse struct {
	JID  string `json:"JID"`
	Name string `json:"Name,omitempty"`
}
