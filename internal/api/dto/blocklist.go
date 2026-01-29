package dto

// BlocklistResponse lista de contatos bloqueados
type BlocklistResponse struct {
	JIDs []string `json:"JIDs" example:"5511999999999@s.whatsapp.net,5521888888888@s.whatsapp.net"`
}

// BlockRequest request para bloquear/desbloquear contato
type BlockRequest struct {
	Phone string `json:"Phone" example:"5511999999999"`
}

// BlockActionResponse resposta da acao de bloquear/desbloquear
type BlockActionResponse struct {
	Details   string   `json:"Details" example:"Contact blocked"`
	Blocklist []string `json:"Blocklist"`
}
