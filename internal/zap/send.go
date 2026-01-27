package zap

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (m *Manager) SendText(ctx context.Context, name, to, text string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(to)
	if err != nil {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	msg := &waE2E.Message{
		Conversation: proto.String(text),
	}

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) SendImage(ctx context.Context, name, to string, data []byte, caption, mimeType string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(to)
	if jid.IsEmpty() {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	uploaded, err := session.Client.Upload(ctx, data, whatsmeow.MediaImage)
	if err != nil {
		return "", fmt.Errorf("failed to upload: %w", err)
	}

	msg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			Caption:       proto.String(caption),
		},
	}

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) SendDocument(ctx context.Context, name, to string, data []byte, filename, mimeType string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(to)
	if jid.IsEmpty() {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	uploaded, err := session.Client.Upload(ctx, data, whatsmeow.MediaDocument)
	if err != nil {
		return "", fmt.Errorf("failed to upload: %w", err)
	}

	msg := &waE2E.Message{
		DocumentMessage: &waE2E.DocumentMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			FileName:      proto.String(filename),
		},
	}

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) SendAudio(ctx context.Context, name, to string, data []byte, mimeType string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(to)
	if jid.IsEmpty() {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	uploaded, err := session.Client.Upload(ctx, data, whatsmeow.MediaAudio)
	if err != nil {
		return "", fmt.Errorf("failed to upload: %w", err)
	}

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
		},
	}

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) SendLocation(ctx context.Context, name, to string, lat, lng float64, locName, address string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(to)
	if jid.IsEmpty() {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(lat),
			DegreesLongitude: proto.Float64(lng),
			Name:             proto.String(locName),
			Address:          proto.String(address),
		},
	}

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) CheckPhone(ctx context.Context, name string, phones []string) ([]types.IsOnWhatsAppResponse, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	return session.Client.IsOnWhatsApp(ctx, phones)
}

func (m *Manager) GetUserInfo(ctx context.Context, name string, jids []types.JID) (map[types.JID]types.UserInfo, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	return session.Client.GetUserInfo(ctx, jids)
}

func (m *Manager) GetProfilePicture(ctx context.Context, name string, jid types.JID) (*types.ProfilePictureInfo, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	return session.Client.GetProfilePictureInfo(ctx, jid, nil)
}

func (m *Manager) SendVideo(ctx context.Context, name, to string, data []byte, caption, mimeType string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(to)
	if jid.IsEmpty() {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	uploaded, err := session.Client.Upload(ctx, data, whatsmeow.MediaVideo)
	if err != nil {
		return "", fmt.Errorf("failed to upload: %w", err)
	}

	msg := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			Caption:       proto.String(caption),
		},
	}

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) SendSticker(ctx context.Context, name, to string, data []byte, mimeType string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(to)
	if jid.IsEmpty() {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	uploaded, err := session.Client.Upload(ctx, data, whatsmeow.MediaImage)
	if err != nil {
		return "", fmt.Errorf("failed to upload: %w", err)
	}

	msg := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
		},
	}

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) SendContact(ctx context.Context, name, to, contactName, vcard string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(to)
	if jid.IsEmpty() {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	msg := &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: proto.String(contactName),
			Vcard:       proto.String(vcard),
		},
	}

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) SendPoll(ctx context.Context, name, to, question string, options []string, multiSelect bool) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(to)
	if jid.IsEmpty() {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	selectCount := 1
	if multiSelect {
		selectCount = len(options)
	}

	msg := session.Client.BuildPollCreation(question, options, selectCount)

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) SendReaction(ctx context.Context, name, to, messageId, emoji string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(to)
	if jid.IsEmpty() {
		jid = types.NewJID(to, types.DefaultUserServer)
	}

	msg := session.Client.BuildReaction(jid, session.Client.Store.ID.ToNonAD(), messageId, emoji)

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) EditMessage(ctx context.Context, name, chatJid, messageId, newText string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(chatJid)
	if jid.IsEmpty() {
		jid = types.NewJID(chatJid, types.DefaultUserServer)
	}

	newContent := &waE2E.Message{
		Conversation: proto.String(newText),
	}

	msg := session.Client.BuildEdit(jid, messageId, newContent)

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to edit: %w", err)
	}

	return resp.ID, nil
}

func (m *Manager) RevokeMessage(ctx context.Context, name, chatJid, messageId string) (string, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(chatJid)
	if jid.IsEmpty() {
		jid = types.NewJID(chatJid, types.DefaultUserServer)
	}

	msg := session.Client.BuildRevoke(jid, session.Client.Store.ID.ToNonAD(), messageId)

	resp, err := session.Client.SendMessage(ctx, jid, msg)
	if err != nil {
		return "", fmt.Errorf("failed to revoke: %w", err)
	}

	return resp.ID, nil
}

// SendPresence sets global presence (available/unavailable)
func (m *Manager) SendPresence(ctx context.Context, name string, presenceType types.Presence) error {
	session, err := m.GetSession(name)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	return session.Client.SendPresence(ctx, presenceType)
}

// SubscribePresence subscribes to presence updates for a contact
func (m *Manager) SubscribePresence(ctx context.Context, name, phone string) error {
	session, err := m.GetSession(name)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(phone)
	if jid.IsEmpty() {
		jid = types.NewJID(phone, types.DefaultUserServer)
	}

	return session.Client.SubscribePresence(ctx, jid)
}

// GetBlocklist returns the list of blocked contacts
func (m *Manager) GetBlocklist(ctx context.Context, name string) (*types.Blocklist, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	return session.Client.GetBlocklist(ctx)
}

// BlockContact blocks a contact
func (m *Manager) BlockContact(ctx context.Context, name, phone string) (*types.Blocklist, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(phone)
	if jid.IsEmpty() {
		jid = types.NewJID(phone, types.DefaultUserServer)
	}

	return session.Client.UpdateBlocklist(ctx, jid, "block")
}

// UnblockContact unblocks a contact
func (m *Manager) UnblockContact(ctx context.Context, name, phone string) (*types.Blocklist, error) {
	session, err := m.GetSession(name)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	jid, _ := types.ParseJID(phone)
	if jid.IsEmpty() {
		jid = types.NewJID(phone, types.DefaultUserServer)
	}

	return session.Client.UpdateBlocklist(ctx, jid, "unblock")
}
