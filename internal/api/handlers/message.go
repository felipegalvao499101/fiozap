package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/fiozap/fiozap/internal/api/dto"
	"github.com/fiozap/fiozap/internal/zap"
	"github.com/go-chi/chi/v5"
	"go.mau.fi/whatsmeow/types"
)

type MessageHandler struct {
	manager *zap.Manager
}

func NewMessageHandler(manager *zap.Manager) *MessageHandler {
	return &MessageHandler{manager: manager}
}

func (h *MessageHandler) SendText(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	msgID, err := h.manager.SendText(r.Context(), name, req.To, req.Text)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageID: msgID})
}

func (h *MessageHandler) SendImage(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	data, err := base64.StdEncoding.DecodeString(req.Image)
	if err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid base64")
		return
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	msgID, err := h.manager.SendImage(r.Context(), name, req.To, data, req.Caption, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageID: msgID})
}

func (h *MessageHandler) SendDocument(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	data, err := base64.StdEncoding.DecodeString(req.Document)
	if err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid base64")
		return
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	msgID, err := h.manager.SendDocument(r.Context(), name, req.To, data, req.Filename, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageID: msgID})
}

func (h *MessageHandler) SendAudio(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendAudioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	data, err := base64.StdEncoding.DecodeString(req.Audio)
	if err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid base64")
		return
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "audio/ogg; codecs=opus"
	}

	msgID, err := h.manager.SendAudio(r.Context(), name, req.To, data, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageID: msgID})
}

func (h *MessageHandler) SendLocation(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	msgID, err := h.manager.SendLocation(r.Context(), name, req.To, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageID: msgID})
}

func (h *MessageHandler) CheckPhone(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.CheckPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	results, err := h.manager.CheckPhone(r.Context(), name, req.Phones)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	list := make([]dto.CheckPhoneResponse, 0, len(results))
	for _, res := range results {
		resp := dto.CheckPhoneResponse{Phone: res.Query, IsOnWhatsApp: res.IsIn}
		if res.IsIn {
			resp.JID = res.JID.String()
		}
		list = append(list, resp)
	}

	dto.Success(w, list)
}

func (h *MessageHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	phone := chi.URLParam(r, "phone")

	jid := types.NewJID(phone, types.DefaultUserServer)
	info, err := h.manager.GetUserInfo(r.Context(), name, []types.JID{jid})
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	userInfo, ok := info[jid]
	if !ok {
		dto.Error(w, http.StatusNotFound, "user not found")
		return
	}

	dto.Success(w, dto.UserInfoResponse{
		JID:    jid.String(),
		Name:   userInfo.VerifiedName.Details.GetVerifiedName(),
		Status: userInfo.Status,
	})
}

func (h *MessageHandler) GetUserAvatar(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	phone := chi.URLParam(r, "phone")

	jid := types.NewJID(phone, types.DefaultUserServer)
	pic, err := h.manager.GetProfilePicture(r.Context(), name, jid)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if pic == nil {
		dto.Error(w, http.StatusNotFound, "no profile picture")
		return
	}

	dto.Success(w, dto.UserInfoResponse{
		JID:        jid.String(),
		PictureURL: pic.URL,
	})
}
