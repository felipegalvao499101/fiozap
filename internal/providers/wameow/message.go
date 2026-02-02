package wameow

import (
	"context"
	"fmt"

	"fiozap/internal/core"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

// SendText envia mensagem de texto
func (m *Manager) SendText(ctx context.Context, session, to, text string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendMessage(ctx, parseJID(to), &waE2E.Message{
		Conversation: proto.String(text),
	})
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// SendImage envia imagem
func (m *Manager) SendImage(ctx context.Context, session, to string, data []byte, caption, mimeType string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	uploaded, err := client.Upload(ctx, data, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	resp, err := client.SendMessage(ctx, parseJID(to), &waE2E.Message{
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
	})
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// SendVideo envia video
func (m *Manager) SendVideo(ctx context.Context, session, to string, data []byte, caption, mimeType string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	uploaded, err := client.Upload(ctx, data, whatsmeow.MediaVideo)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	resp, err := client.SendMessage(ctx, parseJID(to), &waE2E.Message{
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
	})
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// SendAudio envia audio
func (m *Manager) SendAudio(ctx context.Context, session, to string, data []byte, mimeType string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	uploaded, err := client.Upload(ctx, data, whatsmeow.MediaAudio)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	resp, err := client.SendMessage(ctx, parseJID(to), &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// SendDocument envia documento
func (m *Manager) SendDocument(ctx context.Context, session, to string, data []byte, filename, mimeType string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	uploaded, err := client.Upload(ctx, data, whatsmeow.MediaDocument)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	resp, err := client.SendMessage(ctx, parseJID(to), &waE2E.Message{
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
	})
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// SendSticker envia sticker
func (m *Manager) SendSticker(ctx context.Context, session, to string, data []byte, mimeType string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	uploaded, err := client.Upload(ctx, data, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	resp, err := client.SendMessage(ctx, parseJID(to), &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(mimeType),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// SendLocation envia localizacao
func (m *Manager) SendLocation(ctx context.Context, session, to string, lat, lng float64, name, address string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendMessage(ctx, parseJID(to), &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(lat),
			DegreesLongitude: proto.Float64(lng),
			Name:             proto.String(name),
			Address:          proto.String(address),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// SendContact envia contato
func (m *Manager) SendContact(ctx context.Context, session, to, name, vcard string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendMessage(ctx, parseJID(to), &waE2E.Message{
		ContactMessage: &waE2E.ContactMessage{
			DisplayName: proto.String(name),
			Vcard:       proto.String(vcard),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// SendPoll envia enquete
func (m *Manager) SendPoll(ctx context.Context, session, to, question string, options []string, multiSelect bool) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	selectCount := 1
	if multiSelect {
		selectCount = len(options)
	}

	msg := client.BuildPollCreation(question, options, selectCount)
	resp, err := client.SendMessage(ctx, parseJID(to), msg)
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// SendReaction envia reacao
func (m *Manager) SendReaction(ctx context.Context, session, to, messageID, emoji string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	jid := parseJID(to)
	msg := client.BuildReaction(jid, client.Store.ID.ToNonAD(), messageID, emoji)
	resp, err := client.SendMessage(ctx, jid, msg)
	if err != nil {
		return nil, fmt.Errorf("send failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// EditMessage edita mensagem
func (m *Manager) EditMessage(ctx context.Context, session, chat, messageID, newText string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	jid := parseJID(chat)
	msg := client.BuildEdit(jid, messageID, &waE2E.Message{Conversation: proto.String(newText)})
	resp, err := client.SendMessage(ctx, jid, msg)
	if err != nil {
		return nil, fmt.Errorf("edit failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}

// RevokeMessage revoga/apaga mensagem
func (m *Manager) RevokeMessage(ctx context.Context, session, chat, messageID string) (*core.MessageResponse, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	jid := parseJID(chat)
	msg := client.BuildRevoke(jid, client.Store.ID.ToNonAD(), messageID)
	resp, err := client.SendMessage(ctx, jid, msg)
	if err != nil {
		return nil, fmt.Errorf("revoke failed: %w", err)
	}

	return &core.MessageResponse{ID: resp.ID, Timestamp: resp.Timestamp}, nil
}
