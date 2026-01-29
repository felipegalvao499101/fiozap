package domain

import (
	"context"
	"time"
)

// Provider interface para provedores de mensageria (whatsmeow, cloudapi, etc)
type Provider interface {
	// Session management
	CreateSession(ctx context.Context, name string) (Session, error)
	GetSession(name string) (Session, error)
	ListSessions() []Session
	DeleteSession(ctx context.Context, name string) error
	Connect(ctx context.Context, name string) (Session, error)
	Disconnect(name string) error
	Logout(ctx context.Context, name string) error

	// Messages
	SendText(ctx context.Context, session, to, text string) (*MessageResponse, error)
	SendImage(ctx context.Context, session, to string, data []byte, caption, mimeType string) (*MessageResponse, error)
	SendVideo(ctx context.Context, session, to string, data []byte, caption, mimeType string) (*MessageResponse, error)
	SendAudio(ctx context.Context, session, to string, data []byte, mimeType string) (*MessageResponse, error)
	SendDocument(ctx context.Context, session, to string, data []byte, filename, mimeType string) (*MessageResponse, error)
	SendSticker(ctx context.Context, session, to string, data []byte, mimeType string) (*MessageResponse, error)
	SendLocation(ctx context.Context, session, to string, lat, lng float64, name, address string) (*MessageResponse, error)
	SendContact(ctx context.Context, session, to, name, vcard string) (*MessageResponse, error)
	SendPoll(ctx context.Context, session, to, question string, options []string, multiSelect bool) (*MessageResponse, error)
	SendReaction(ctx context.Context, session, to, messageID, emoji string) (*MessageResponse, error)
	EditMessage(ctx context.Context, session, chat, messageID, newText string) (*MessageResponse, error)
	RevokeMessage(ctx context.Context, session, chat, messageID string) (*MessageResponse, error)

	// Chat
	MarkRead(ctx context.Context, session, chatJID string, messageIDs []string) error
	SendTyping(ctx context.Context, session, chatJID string, composing bool) error
	SendRecording(ctx context.Context, session, chatJID string, recording bool) error
	SetDisappearingTimer(ctx context.Context, session, chatJID string, duration time.Duration) error
	SendPresence(ctx context.Context, session string, available bool) error
	SubscribePresence(ctx context.Context, session, phone string) error
	RejectCall(ctx context.Context, session, callFrom, callID string) error

	// Groups
	CreateGroup(ctx context.Context, session, name string, participants []string) (*GroupInfo, error)
	GetGroups(ctx context.Context, session string) ([]*GroupInfo, error)
	GetGroupInfo(ctx context.Context, session, groupJID string) (*GroupInfo, error)
	SetGroupName(ctx context.Context, session, groupJID, name string) error
	SetGroupTopic(ctx context.Context, session, groupJID, topic string) error
	SetGroupPhoto(ctx context.Context, session, groupJID, photoBase64 string) (string, error)
	LeaveGroup(ctx context.Context, session, groupJID string) error
	GetGroupInviteLink(ctx context.Context, session, groupJID string, reset bool) (string, error)
	JoinGroupWithLink(ctx context.Context, session, code string) (string, error)
	GetGroupInfoFromLink(ctx context.Context, session, code string) (*GroupInfo, error)
	AddParticipants(ctx context.Context, session, groupJID string, phones []string) error
	RemoveParticipants(ctx context.Context, session, groupJID string, phones []string) error
	PromoteParticipants(ctx context.Context, session, groupJID string, phones []string) error
	DemoteParticipants(ctx context.Context, session, groupJID string, phones []string) error
	SetGroupAnnounce(ctx context.Context, session, groupJID string, announce bool) error
	SetGroupLocked(ctx context.Context, session, groupJID string, locked bool) error

	// Users
	CheckPhone(ctx context.Context, session string, phones []string) ([]PhoneCheck, error)
	GetUserInfo(ctx context.Context, session string, phones []string) (map[string]UserInfo, error)
	GetProfilePicture(ctx context.Context, session, phone string) (*ProfilePicture, error)
	GetBlocklist(ctx context.Context, session string) ([]string, error)
	BlockContact(ctx context.Context, session, phone string) ([]string, error)
	UnblockContact(ctx context.Context, session, phone string) ([]string, error)
}
