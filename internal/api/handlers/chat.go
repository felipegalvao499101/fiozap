package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fiozap/fiozap/internal/api/dto"
	"github.com/fiozap/fiozap/internal/zap"
	"github.com/go-chi/chi/v5"
	"go.mau.fi/whatsmeow/types"
)

type ChatHandler struct {
	manager *zap.Manager
}

func NewChatHandler(manager *zap.Manager) *ChatHandler {
	return &ChatHandler{manager: manager}
}

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

	if err := h.manager.MarkRead(r.Context(), name, chatJid, req.Id); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Messages marked as read"})
}

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

	session, err := h.manager.GetSession(name)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if session.Client == nil || !session.Client.IsConnected() {
		dto.Error(w, http.StatusInternalServerError, "session not connected")
		return
	}

	jid, _ := types.ParseJID(req.Phone)
	if jid.IsEmpty() {
		jid = types.NewJID(req.Phone, types.DefaultUserServer)
	}

	var media types.ChatPresenceMedia
	if req.Media == "audio" {
		media = types.ChatPresenceMediaAudio
	} else {
		media = types.ChatPresenceMediaText
	}

	err = session.Client.SendChatPresence(r.Context(), jid, types.ChatPresence(req.State), media)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Chat presence set successfully"})
}

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

	if err := h.manager.SetDisappearingTimer(r.Context(), name, chatJid, duration); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Disappearing timer set successfully"})
}

// Presence handlers (global)
func (h *ChatHandler) SendPresence(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.PresenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	var presence types.Presence
	switch req.Type {
	case "available":
		presence = types.PresenceAvailable
	case "unavailable":
		presence = types.PresenceUnavailable
	default:
		dto.Error(w, http.StatusBadRequest, "invalid presence Type. Allowed: available, unavailable")
		return
	}

	if err := h.manager.SendPresence(r.Context(), name, presence); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Presence set successfully"})
}

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

	if err := h.manager.SubscribePresence(r.Context(), name, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Subscribed to presence updates"})
}

// Blocklist handlers
func (h *ChatHandler) GetBlocklist(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	blocklist, err := h.manager.GetBlocklist(r.Context(), name)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	jids := make([]string, len(blocklist.JIDs))
	for i, jid := range blocklist.JIDs {
		jids[i] = jid.String()
	}

	dto.Success(w, dto.BlocklistResponse{JIDs: jids})
}

func (h *ChatHandler) BlockContact(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.BlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	blocklist, err := h.manager.BlockContact(r.Context(), name, req.Phone)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	jids := make([]string, len(blocklist.JIDs))
	for i, jid := range blocklist.JIDs {
		jids[i] = jid.String()
	}

	dto.Success(w, map[string]interface{}{"Details": "Contact blocked", "Blocklist": jids})
}

func (h *ChatHandler) UnblockContact(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.BlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	blocklist, err := h.manager.UnblockContact(r.Context(), name, req.Phone)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	jids := make([]string, len(blocklist.JIDs))
	for i, jid := range blocklist.JIDs {
		jids[i] = jid.String()
	}

	dto.Success(w, map[string]interface{}{"Details": "Contact unblocked", "Blocklist": jids})
}
