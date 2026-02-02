package wameow

import (
	"context"
	"encoding/base64"
	"fmt"

	"fiozap/internal/core"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// CreateGroup cria um grupo
func (m *Manager) CreateGroup(ctx context.Context, session, name string, participants []string) (*core.GroupInfo, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	jids := make([]types.JID, len(participants))
	for i, phone := range participants {
		jids[i] = parseJID(phone)
	}

	info, err := client.CreateGroup(ctx, whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: jids,
	})
	if err != nil {
		return nil, fmt.Errorf("create group failed: %w", err)
	}

	return groupToInfo(info), nil
}

// GetGroups lista grupos
func (m *Manager) GetGroups(ctx context.Context, session string) ([]*core.GroupInfo, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	groups, err := client.GetJoinedGroups(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*core.GroupInfo, len(groups))
	for i, g := range groups {
		result[i] = groupToInfo(g)
	}
	return result, nil
}

// GetGroupInfo retorna info do grupo
func (m *Manager) GetGroupInfo(ctx context.Context, session, groupJID string) (*core.GroupInfo, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return nil, fmt.Errorf("invalid group JID: %w", err)
	}

	info, err := client.GetGroupInfo(ctx, jid)
	if err != nil {
		return nil, err
	}

	return groupToInfo(info), nil
}

// SetGroupName altera nome do grupo
func (m *Manager) SetGroupName(ctx context.Context, session, groupJID, name string) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	return client.SetGroupName(ctx, jid, name)
}

// SetGroupTopic altera descricao do grupo
func (m *Manager) SetGroupTopic(ctx context.Context, session, groupJID, topic string) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	return client.SetGroupDescription(ctx, jid, topic)
}

// SetGroupPhoto altera foto do grupo
func (m *Manager) SetGroupPhoto(ctx context.Context, session, groupJID, photoBase64 string) (string, error) {
	client, err := m.getClient(session)
	if err != nil {
		return "", err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return "", fmt.Errorf("invalid group JID: %w", err)
	}

	data, err := base64.StdEncoding.DecodeString(photoBase64)
	if err != nil {
		return "", fmt.Errorf("invalid base64: %w", err)
	}

	return client.SetGroupPhoto(ctx, jid, data)
}

// LeaveGroup sai do grupo
func (m *Manager) LeaveGroup(ctx context.Context, session, groupJID string) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	return client.LeaveGroup(ctx, jid)
}

// GetGroupInviteLink retorna link de convite
func (m *Manager) GetGroupInviteLink(ctx context.Context, session, groupJID string, reset bool) (string, error) {
	client, err := m.getClient(session)
	if err != nil {
		return "", err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return "", fmt.Errorf("invalid group JID: %w", err)
	}

	return client.GetGroupInviteLink(ctx, jid, reset)
}

// JoinGroupWithLink entra no grupo pelo link
func (m *Manager) JoinGroupWithLink(ctx context.Context, session, code string) (string, error) {
	client, err := m.getClient(session)
	if err != nil {
		return "", err
	}

	jid, err := client.JoinGroupWithLink(ctx, code)
	if err != nil {
		return "", err
	}

	return jid.String(), nil
}

// GetGroupInfoFromLink retorna info do grupo pelo link
func (m *Manager) GetGroupInfoFromLink(ctx context.Context, session, code string) (*core.GroupInfo, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	info, err := client.GetGroupInfoFromLink(ctx, code)
	if err != nil {
		return nil, err
	}

	return groupToInfo(info), nil
}

// AddParticipants adiciona participantes
func (m *Manager) AddParticipants(ctx context.Context, session, groupJID string, phones []string) error {
	return m.updateParticipants(ctx, session, groupJID, phones, whatsmeow.ParticipantChangeAdd)
}

// RemoveParticipants remove participantes
func (m *Manager) RemoveParticipants(ctx context.Context, session, groupJID string, phones []string) error {
	return m.updateParticipants(ctx, session, groupJID, phones, whatsmeow.ParticipantChangeRemove)
}

// PromoteParticipants promove a admin
func (m *Manager) PromoteParticipants(ctx context.Context, session, groupJID string, phones []string) error {
	return m.updateParticipants(ctx, session, groupJID, phones, whatsmeow.ParticipantChangePromote)
}

// DemoteParticipants rebaixa admin
func (m *Manager) DemoteParticipants(ctx context.Context, session, groupJID string, phones []string) error {
	return m.updateParticipants(ctx, session, groupJID, phones, whatsmeow.ParticipantChangeDemote)
}

func (m *Manager) updateParticipants(ctx context.Context, session, groupJID string, phones []string, action whatsmeow.ParticipantChange) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	jids := make([]types.JID, len(phones))
	for i, phone := range phones {
		jids[i] = parseJID(phone)
	}

	_, err = client.UpdateGroupParticipants(ctx, jid, jids, action)
	return err
}

// SetGroupAnnounce define se so admins podem enviar mensagens
func (m *Manager) SetGroupAnnounce(ctx context.Context, session, groupJID string, announce bool) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	return client.SetGroupAnnounce(ctx, jid, announce)
}

// SetGroupLocked define se so admins podem editar info
func (m *Manager) SetGroupLocked(ctx context.Context, session, groupJID string, locked bool) error {
	client, err := m.getClient(session)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(groupJID)
	if err != nil {
		return fmt.Errorf("invalid group JID: %w", err)
	}

	return client.SetGroupLocked(ctx, jid, locked)
}

func groupToInfo(g *types.GroupInfo) *core.GroupInfo {
	participants := make([]core.GroupParticipant, len(g.Participants))
	for i, p := range g.Participants {
		participants[i] = core.GroupParticipant{
			JID:          p.JID.String(),
			IsAdmin:      p.IsAdmin,
			IsSuperAdmin: p.IsSuperAdmin,
		}
	}

	return &core.GroupInfo{
		JID:          g.JID.String(),
		Name:         g.Name,
		Topic:        g.Topic,
		GroupCreated: g.GroupCreated,
		OwnerJID:     g.OwnerJID.String(),
		Participants: participants,
		IsAnnounce:   g.IsAnnounce,
		IsLocked:     g.IsLocked,
	}
}
