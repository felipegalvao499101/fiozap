package handlers

import (
	"net/http"

	"fiozap/internal/api/dto"
	"fiozap/internal/domain"
)

type ProfileHandler struct {
	provider domain.Provider
}

func NewProfileHandler(provider domain.Provider) *ProfileHandler {
	return &ProfileHandler{provider: provider}
}

// GetContactQRLink godoc
// @Summary      Obter QR link de contato
// @Description  Retorna link QR para adicionar contato (nao implementado)
// @Tags         profile
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/profile/qrlink [get]
func (h *ProfileHandler) GetContactQRLink(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// ResolveContactQRLink godoc
// @Summary      Resolver QR link
// @Description  Resolve um link QR de contato (nao implementado)
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.ResolveContactQRRequest true "Link QR"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/profile/qrlink/resolve [post]
func (h *ProfileHandler) ResolveContactQRLink(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// SetStatusMessage godoc
// @Summary      Alterar recado
// @Description  Altera o recado/status do perfil (nao implementado)
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SetStatusMessageRequest true "Novo status"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/profile/status [put]
func (h *ProfileHandler) SetStatusMessage(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// ResolveBusinessMessageLink godoc
// @Summary      Resolver link comercial
// @Description  Resolve um link de mensagem comercial (nao implementado)
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.ResolveBusinessLinkRequest true "Link"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/profile/business/resolve [post]
func (h *ProfileHandler) ResolveBusinessMessageLink(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// GetBusinessProfile godoc
// @Summary      Perfil comercial
// @Description  Retorna perfil comercial de um contato (nao implementado)
// @Tags         profile
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        phone path string true "Numero do telefone"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/users/{phone}/business [get]
func (h *ProfileHandler) GetBusinessProfile(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}
