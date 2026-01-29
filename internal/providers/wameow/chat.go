package wameow

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow/types"
)

// MarkRead marca mensagens como lidas
func (m *Manager) MarkRead(ctx context.Context, session, chatJID string, messageIDs []string) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	jid := parseJID(chatJID)
	ids := make([]types.MessageID, len(messageIDs))
	for i, id := range messageIDs {
		ids[i] = types.MessageID(id)
	}

	return client.MarkRead(ctx, ids, time.Now(), jid, jid)
}

// SendTyping envia indicador de digitacao
func (m *Manager) SendTyping(ctx context.Context, session, chatJID string, composing bool) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	state := types.ChatPresencePaused
	if composing {
		state = types.ChatPresenceComposing
	}

	return client.SendChatPresence(ctx, parseJID(chatJID), state, types.ChatPresenceMediaText)
}

// SendRecording envia indicador de gravacao de audio
func (m *Manager) SendRecording(ctx context.Context, session, chatJID string, recording bool) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	state := types.ChatPresencePaused
	if recording {
		state = types.ChatPresenceComposing
	}

	return client.SendChatPresence(ctx, parseJID(chatJID), state, types.ChatPresenceMediaAudio)
}

// SetDisappearingTimer define timer de mensagens temporarias
func (m *Manager) SetDisappearingTimer(ctx context.Context, session, chatJID string, duration time.Duration) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	return client.SetDisappearingTimer(ctx, parseJID(chatJID), duration, time.Now())
}

// SendPresence envia presenca online/offline
func (m *Manager) SendPresence(ctx context.Context, session string, available bool) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	presence := types.PresenceUnavailable
	if available {
		presence = types.PresenceAvailable
	}

	return client.SendPresence(ctx, presence)
}

// SubscribePresence se inscreve para receber presenca de um contato
func (m *Manager) SubscribePresence(ctx context.Context, session, phone string) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	return client.SubscribePresence(ctx, parseJID(phone))
}

// RejectCall rejeita chamada
func (m *Manager) RejectCall(ctx context.Context, session, callFrom, callID string) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(callFrom)
	if err != nil {
		return fmt.Errorf("invalid JID: %w", err)
	}

	return client.RejectCall(ctx, jid, callID)
}
