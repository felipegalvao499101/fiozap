package zap

import (
	"context"
	"encoding/base64"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type GroupInfo struct {
	JID          string
	Name         string
	Topic        string
	GroupCreated string
	OwnerJID     string
	Participants []GroupParticipant
	IsAnnounce   bool
	IsLocked     bool
}

type GroupParticipant struct {
	JID          string
	IsAdmin      bool
	IsSuperAdmin bool
}

func (m *Manager) CreateGroup(ctx context.Context, name string, sessionName string, participants []string) (*GroupInfo, error) {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	jids := make([]types.JID, len(participants))
	for i, p := range participants {
		jids[i] = types.NewJID(p, types.DefaultUserServer)
	}

	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: jids,
	}

	info, err := session.Client.CreateGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return groupInfoToStruct(info), nil
}

func (m *Manager) GetGroups(ctx context.Context, sessionName string) ([]*GroupInfo, error) {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	groups, err := session.Client.GetJoinedGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	result := make([]*GroupInfo, len(groups))
	for i, g := range groups {
		result[i] = groupInfoToStruct(g)
	}

	return result, nil
}

func (m *Manager) GetGroupInfo(ctx context.Context, sessionName, groupJid string) (*GroupInfo, error) {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return nil, fmt.Errorf("invalid group jid: %w", err)
	}

	info, err := session.Client.GetGroupInfo(ctx, jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}

	return groupInfoToStruct(info), nil
}

func (m *Manager) SetGroupName(ctx context.Context, sessionName, groupJid, name string) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return fmt.Errorf("invalid group jid: %w", err)
	}

	return session.Client.SetGroupName(ctx, jid, name)
}

func (m *Manager) SetGroupDescription(ctx context.Context, sessionName, groupJid, description string) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return fmt.Errorf("invalid group jid: %w", err)
	}

	return session.Client.SetGroupDescription(ctx, jid, description)
}

func (m *Manager) SetGroupPhoto(ctx context.Context, sessionName, groupJid, photoBase64 string) (string, error) {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return "", fmt.Errorf("invalid group jid: %w", err)
	}

	data, err := base64.StdEncoding.DecodeString(photoBase64)
	if err != nil {
		return "", fmt.Errorf("invalid base64: %w", err)
	}

	return session.Client.SetGroupPhoto(ctx, jid, data)
}

func (m *Manager) LeaveGroup(ctx context.Context, sessionName, groupJid string) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return fmt.Errorf("invalid group jid: %w", err)
	}

	return session.Client.LeaveGroup(ctx, jid)
}

func (m *Manager) GetGroupInviteLink(ctx context.Context, sessionName, groupJid string, reset bool) (string, error) {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return "", fmt.Errorf("invalid group jid: %w", err)
	}

	return session.Client.GetGroupInviteLink(ctx, jid, reset)
}

func (m *Manager) JoinGroupWithLink(ctx context.Context, sessionName, code string) (string, error) {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return "", err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return "", fmt.Errorf("session not connected")
	}

	jid, err := session.Client.JoinGroupWithLink(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to join group: %w", err)
	}

	return jid.String(), nil
}

func (m *Manager) GetGroupInfoFromLink(ctx context.Context, sessionName, code string) (*GroupInfo, error) {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return nil, err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}

	info, err := session.Client.GetGroupInfoFromLink(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}

	return groupInfoToStruct(info), nil
}

func (m *Manager) AddGroupParticipants(ctx context.Context, sessionName, groupJid string, participants []string) error {
	return m.updateGroupParticipants(ctx, sessionName, groupJid, participants, whatsmeow.ParticipantChangeAdd)
}

func (m *Manager) RemoveGroupParticipants(ctx context.Context, sessionName, groupJid string, participants []string) error {
	return m.updateGroupParticipants(ctx, sessionName, groupJid, participants, whatsmeow.ParticipantChangeRemove)
}

func (m *Manager) PromoteGroupParticipants(ctx context.Context, sessionName, groupJid string, participants []string) error {
	return m.updateGroupParticipants(ctx, sessionName, groupJid, participants, whatsmeow.ParticipantChangePromote)
}

func (m *Manager) DemoteGroupParticipants(ctx context.Context, sessionName, groupJid string, participants []string) error {
	return m.updateGroupParticipants(ctx, sessionName, groupJid, participants, whatsmeow.ParticipantChangeDemote)
}

func (m *Manager) updateGroupParticipants(ctx context.Context, sessionName, groupJid string, participants []string, action whatsmeow.ParticipantChange) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return fmt.Errorf("invalid group jid: %w", err)
	}

	jids := make([]types.JID, len(participants))
	for i, p := range participants {
		jids[i] = types.NewJID(p, types.DefaultUserServer)
	}

	_, err = session.Client.UpdateGroupParticipants(ctx, jid, jids, action)
	return err
}

func (m *Manager) SetGroupAnnounce(ctx context.Context, sessionName, groupJid string, announce bool) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return fmt.Errorf("invalid group jid: %w", err)
	}

	return session.Client.SetGroupAnnounce(ctx, jid, announce)
}

func (m *Manager) SetGroupLocked(ctx context.Context, sessionName, groupJid string, locked bool) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return fmt.Errorf("invalid group jid: %w", err)
	}

	return session.Client.SetGroupLocked(ctx, jid, locked)
}

func (m *Manager) SetGroupJoinApproval(ctx context.Context, sessionName, groupJid string, approval bool) error {
	session, err := m.GetSession(sessionName)
	if err != nil {
		return err
	}

	if session.Client == nil || !session.Client.IsConnected() {
		return fmt.Errorf("session not connected")
	}

	jid, err := types.ParseJID(groupJid)
	if err != nil {
		return fmt.Errorf("invalid group jid: %w", err)
	}

	return session.Client.SetGroupJoinApprovalMode(ctx, jid, approval)
}

func groupInfoToStruct(info *types.GroupInfo) *GroupInfo {
	participants := make([]GroupParticipant, len(info.Participants))
	for i, p := range info.Participants {
		participants[i] = GroupParticipant{
			JID:          p.JID.String(),
			IsAdmin:      p.IsAdmin,
			IsSuperAdmin: p.IsSuperAdmin,
		}
	}

	return &GroupInfo{
		JID:          info.JID.String(),
		Name:         info.Name,
		Topic:        info.Topic,
		GroupCreated: info.GroupCreated.Format("2006-01-02T15:04:05-07:00"),
		OwnerJID:     info.OwnerJID.String(),
		Participants: participants,
		IsAnnounce:   info.IsAnnounce,
		IsLocked:     info.IsLocked,
	}
}
