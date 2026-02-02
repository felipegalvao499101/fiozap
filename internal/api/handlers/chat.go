package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"fiozap/internal/api/dto"
	"fiozap/internal/core"

	"github.com/go-chi/chi/v5"
)

type ChatHandler struct {
	provider core.Provider
}

func NewChatHandler(provider core.Provider) *ChatHandler {
	return &ChatHandler{provider: provider}
}

// MarkRead godoc
// @Summary      Marcar como lida
// @Description  Marca mensagens como lidas
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.MarkReadRequest true "IDs das mensagens"
// @Success      200 {object} dto.Response{data=dto.ChatActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/chat/markread [post]
func (h *ChatHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.MarkReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if len(req.Id) == 0 {
		dto.Error(w, http.StatusBadRequest, "missing Id in Payload")
		return
	}

	chatJid := req.ChatPhone
	if chatJid == "" {
		dto.Error(w, http.StatusBadRequest, "missing ChatPhone in Payload")
		return
	}

	if err := h.provider.MarkRead(r.Context(), name, chatJid, req.Id); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ChatActionResponse{Details: "Messages marked as read"})
}

// Presence godoc
// @Summary      Enviar presenca no chat
// @Description  Envia presenca de digitando/gravando audio para um chat
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.ChatPresenceRequest true "Dados da presenca"
// @Success      200 {object} dto.Response{data=dto.ChatActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/chat/presence [post]
func (h *ChatHandler) Presence(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.ChatPresenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	if req.State == "" {
		dto.Error(w, http.StatusBadRequest, "missing State in Payload")
		return
	}

	composing := req.State == "composing"
	var err error

	if req.Media == "audio" {
		err = h.provider.SendRecording(r.Context(), name, req.Phone, composing)
	} else {
		err = h.provider.SendTyping(r.Context(), name, req.Phone, composing)
	}

	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ChatActionResponse{Details: "Chat presence set successfully"})
}

// SetDisappearing godoc
// @Summary      Configurar mensagens temporarias
// @Description  Define tempo de expiracao das mensagens no chat
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        chatJid path string true "JID do chat"
// @Param        request body dto.DisappearingRequest true "Duracao"
// @Success      200 {object} dto.Response{data=dto.ChatActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/chat/{chatJid}/disappearing [put]
func (h *ChatHandler) SetDisappearing(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	chatJid := chi.URLParam(r, "chatJid")

	var req dto.DisappearingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	var duration time.Duration
	switch req.Duration {
	case "24h":
		duration = 24 * time.Hour
	case "7d":
		duration = 7 * 24 * time.Hour
	case "90d":
		duration = 90 * 24 * time.Hour
	case "off":
		duration = 0
	default:
		dto.Error(w, http.StatusBadRequest, "invalid Duration. Allowed: 24h, 7d, 90d, off")
		return
	}

	if err := h.provider.SetDisappearingTimer(r.Context(), name, chatJid, duration); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ChatActionResponse{Details: "Disappearing timer set successfully"})
}

// SendPresence godoc
// @Summary      Enviar presenca global
// @Description  Define status online/offline da sessao
// @Tags         presence
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.PresenceRequest true "Tipo de presenca"
// @Success      200 {object} dto.Response{data=dto.ChatActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/presence [post]
func (h *ChatHandler) SendPresence(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.PresenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	var available bool
	switch req.Type {
	case "available":
		available = true
	case "unavailable":
		available = false
	default:
		dto.Error(w, http.StatusBadRequest, "invalid presence Type. Allowed: available, unavailable")
		return
	}

	if err := h.provider.SendPresence(r.Context(), name, available); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ChatActionResponse{Details: "Presence set successfully"})
}

// SubscribePresence godoc
// @Summary      Inscrever em presenca
// @Description  Se inscreve para receber atualizacoes de presenca de um contato
// @Tags         presence
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.PresenceSubscribeRequest true "Numero do contato"
// @Success      200 {object} dto.Response{data=dto.ChatActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/presence/subscribe [post]
func (h *ChatHandler) SubscribePresence(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.PresenceSubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	if err := h.provider.SubscribePresence(r.Context(), name, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ChatActionResponse{Details: "Subscribed to presence updates"})
}
