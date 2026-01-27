package zap

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow/types"
)

func (m *Manager) MarkRead(ctx context.Context, sessionName, chatJid string, messageIds []string) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(chatJid)
	if err != nil {
		jid = types.NewJID(chatJid, types.DefaultUserServer)
	}

	ids := make([]types.MessageID, len(messageIds))
	for i, id := range messageIds {
		ids[i] = types.MessageID(id)
	}

	return session.Client.MarkRead(ctx, ids, time.Now(), jid, jid)
}

func (m *Manager) SendTyping(ctx context.Context, sessionName, chatJid string, composing bool) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(chatJid)
	if err != nil {
		jid = types.NewJID(chatJid, types.DefaultUserServer)
	}

	state := types.ChatPresencePaused
	if composing {
		state = types.ChatPresenceComposing
	}

	return session.Client.SendChatPresence(ctx, jid, state, types.ChatPresenceMediaText)
}

func (m *Manager) SendRecording(ctx context.Context, sessionName, chatJid string, composing bool) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(chatJid)
	if err != nil {
		jid = types.NewJID(chatJid, types.DefaultUserServer)
	}

	state := types.ChatPresencePaused
	if composing {
		state = types.ChatPresenceComposing
	}

	return session.Client.SendChatPresence(ctx, jid, state, types.ChatPresenceMediaAudio)
}

func (m *Manager) SetDisappearingTimer(ctx context.Context, sessionName, chatJid string, duration time.Duration) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(chatJid)
	if err != nil {
		jid = types.NewJID(chatJid, types.DefaultUserServer)
	}

	return session.Client.SetDisappearingTimer(ctx, jid, duration, time.Now())
}
