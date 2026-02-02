package handlers

import (
	"net/http"

	"fiozap/internal/api/dto"
	"fiozap/internal/core"
)

type PrivacyHandler struct {
	provider core.Provider
}

func NewPrivacyHandler(provider core.Provider) *PrivacyHandler {
	return &PrivacyHandler{provider: provider}
}

// GetSettings godoc
// @Summary      Obter configuracoes de privacidade
// @Description  Retorna configuracoes de privacidade da sessao (nao implementado)
// @Tags         privacy
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/privacy [get]
func (h *PrivacyHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// SetSetting godoc
// @Summary      Alterar privacidade
// @Description  Altera uma configuracao de privacidade (nao implementado)
// @Tags         privacy
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SetPrivacyRequest true "Configuracao"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/privacy [put]
func (h *PrivacyHandler) SetSetting(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// GetStatusPrivacy godoc
// @Summary      Privacidade do status
// @Description  Retorna configuracoes de privacidade do status (nao implementado)
// @Tags         privacy
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/privacy/status [get]
func (h *PrivacyHandler) GetStatusPrivacy(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}
