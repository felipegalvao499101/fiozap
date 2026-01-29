package wameow

import (
	"context"

	"fiozap/internal/domain"

	"go.mau.fi/whatsmeow/types"
)

// CheckPhone verifica se numeros estao no WhatsApp
func (m *Manager) CheckPhone(ctx context.Context, session string, phones []string) ([]domain.PhoneCheck, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	results, err := client.IsOnWhatsApp(ctx, phones)
	if err != nil {
		return nil, err
	}

	checks := make([]domain.PhoneCheck, len(results))
	for i, r := range results {
		checks[i] = domain.PhoneCheck{
			Phone:        r.Query,
			IsOnWhatsApp: r.IsIn,
			JID:          r.JID.String(),
		}
	}
	return checks, nil
}

// GetUserInfo retorna informacoes de usuarios
func (m *Manager) GetUserInfo(ctx context.Context, session string, phones []string) (map[string]domain.UserInfo, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	jids := make([]types.JID, len(phones))
	for i, phone := range phones {
		jids[i] = parseJID(phone)
	}

	info, err := client.GetUserInfo(ctx, jids)
	if err != nil {
		return nil, err
	}

	result := make(map[string]domain.UserInfo)
	for jid, ui := range info {
		devices := make([]string, len(ui.Devices))
		for i, d := range ui.Devices {
			devices[i] = d.String()
		}
		result[jid.String()] = domain.UserInfo{
			Status:    ui.Status,
			PictureID: ui.PictureID,
			Devices:   devices,
		}
	}
	return result, nil
}

// GetProfilePicture retorna foto de perfil
func (m *Manager) GetProfilePicture(ctx context.Context, session, phone string) (*domain.ProfilePicture, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	pic, err := client.GetProfilePictureInfo(ctx, parseJID(phone), nil)
	if err != nil {
		return nil, err
	}
	if pic == nil {
		return nil, nil
	}

	return &domain.ProfilePicture{URL: pic.URL, ID: pic.ID}, nil
}

// GetBlocklist retorna lista de bloqueados
func (m *Manager) GetBlocklist(ctx context.Context, session string) ([]string, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	blocklist, err := client.GetBlocklist(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(blocklist.JIDs))
	for i, jid := range blocklist.JIDs {
		result[i] = jid.String()
	}
	return result, nil
}

// BlockContact bloqueia contato
func (m *Manager) BlockContact(ctx context.Context, session, phone string) ([]string, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	blocklist, err := client.UpdateBlocklist(ctx, parseJID(phone), "block")
	if err != nil {
		return nil, err
	}

	result := make([]string, len(blocklist.JIDs))
	for i, jid := range blocklist.JIDs {
		result[i] = jid.String()
	}
	return result, nil
}

// UnblockContact desbloqueia contato
func (m *Manager) UnblockContact(ctx context.Context, session, phone string) ([]string, error) {
	client, err := m.getClient(session)
	if err != nil {
		return nil, err
	}

	blocklist, err := client.UpdateBlocklist(ctx, parseJID(phone), "unblock")
	if err != nil {
		return nil, err
	}

	result := make([]string, len(blocklist.JIDs))
	for i, jid := range blocklist.JIDs {
		result[i] = jid.String()
	}
	return result, nil
}
