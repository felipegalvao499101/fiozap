package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fiozap/fiozap/internal/api/dto"
	"github.com/fiozap/fiozap/internal/zap"
	"github.com/go-chi/chi/v5"
	"github.com/skip2/go-qrcode"
)

type SessionHandler struct {
	manager *zap.Manager
}

func NewSessionHandler(manager *zap.Manager) *SessionHandler {
	return &SessionHandler{manager: manager}
}

func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Name == "" {
		dto.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	session, err := h.manager.CreateSession(r.Context(), req.Name)
	if err != nil {
		dto.Error(w, http.StatusConflict, err.Error())
		return
	}

	// Retorna token apenas na criacao
	resp := sessionToDTO(session)
	resp.Token = session.Token
	dto.Created(w, resp)
}

func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	sessions := h.manager.ListSessions()

	list := make([]dto.SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		list = append(list, sessionToDTO(s))
	}

	dto.Success(w, list)
}

func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	session, err := h.manager.GetSession(name)
	if err != nil {
		dto.Error(w, http.StatusNotFound, err.Error())
		return
	}

	dto.Success(w, sessionToDTO(session))
}

func (h *SessionHandler) Connect(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	session, err := h.manager.Connect(r.Context(), name)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, sessionToDTO(session))
}

func (h *SessionHandler) GetQR(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	session, err := h.manager.GetSession(name)
	if err != nil {
		dto.Error(w, http.StatusNotFound, err.Error())
		return
	}

	code := session.GetQRCode()
	if code == "" {
		dto.Error(w, http.StatusNotFound, "no QR code available")
		return
	}

	if r.URL.Query().Get("format") == "image" {
		png, err := qrcode.Encode(code, qrcode.Medium, 256)
		if err != nil {
			dto.Error(w, http.StatusInternalServerError, "failed to generate image")
			return
		}
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(png)
		return
	}

	dto.Success(w, dto.QRResponse{Code: code})
}

func (h *SessionHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.manager.Disconnect(name); err != nil {
		dto.Error(w, http.StatusNotFound, err.Error())
		return
	}

	session, _ := h.manager.GetSession(name)
	dto.Success(w, sessionToDTO(session))
}

func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.manager.DeleteSession(r.Context(), name); err != nil {
		dto.Error(w, http.StatusNotFound, err.Error())
		return
	}

	dto.Success(w, map[string]string{"message": "deleted"})
}

func sessionToDTO(s *zap.Session) dto.SessionResponse {
	return dto.SessionResponse{
		Name:      s.Name,
		JID:       s.GetJID(),
		Phone:     s.GetPhone(),
		PushName:  s.GetPushName(),
		Connected: s.IsConnected(),
	}
}
