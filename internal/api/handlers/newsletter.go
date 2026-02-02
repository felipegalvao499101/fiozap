package handlers

import (
	"net/http"

	"fiozap/internal/api/dto"
	"fiozap/internal/core"
)

type NewsletterHandler struct {
	provider core.Provider
}

func NewNewsletterHandler(provider core.Provider) *NewsletterHandler {
	return &NewsletterHandler{provider: provider}
}

// Create godoc
// @Summary      Criar canal
// @Description  Cria um novo canal/newsletter (nao implementado)
// @Tags         newsletters
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.CreateNewsletterRequest true "Dados do canal"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/newsletters [post]
func (h *NewsletterHandler) Create(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// List godoc
// @Summary      Listar canais
// @Description  Lista canais que a sessao segue (nao implementado)
// @Tags         newsletters
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/newsletters [get]
func (h *NewsletterHandler) List(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// Get godoc
// @Summary      Obter canal
// @Description  Retorna informacoes de um canal (nao implementado)
// @Tags         newsletters
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        newsletterJid path string true "JID do canal"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/newsletters/{newsletterJid} [get]
func (h *NewsletterHandler) Get(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// Follow godoc
// @Summary      Seguir canal
// @Description  Segue um canal (nao implementado)
// @Tags         newsletters
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        newsletterJid path string true "JID do canal"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/newsletters/{newsletterJid}/follow [post]
func (h *NewsletterHandler) Follow(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// Unfollow godoc
// @Summary      Deixar de seguir
// @Description  Deixa de seguir um canal (nao implementado)
// @Tags         newsletters
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        newsletterJid path string true "JID do canal"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/newsletters/{newsletterJid}/unfollow [post]
func (h *NewsletterHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// ToggleMute godoc
// @Summary      Silenciar canal
// @Description  Ativa/desativa notificacoes do canal (nao implementado)
// @Tags         newsletters
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        newsletterJid path string true "JID do canal"
// @Param        request body dto.NewsletterMuteRequest true "Estado mute"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/newsletters/{newsletterJid}/mute [put]
func (h *NewsletterHandler) ToggleMute(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// SendReaction godoc
// @Summary      Reagir em mensagem
// @Description  Envia reacao em mensagem do canal (nao implementado)
// @Tags         newsletters
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        newsletterJid path string true "JID do canal"
// @Param        request body dto.NewsletterReactionRequest true "Dados da reacao"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/newsletters/{newsletterJid}/reaction [post]
func (h *NewsletterHandler) SendReaction(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}
