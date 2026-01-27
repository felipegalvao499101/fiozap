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
