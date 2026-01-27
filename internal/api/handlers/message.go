package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

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
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	msgId, err := h.manager.SendText(r.Context(), name, req.Phone, req.Body)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) SendImage(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	data, err := decodeBase64Data(req.Image)
	if err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode base64 encoded data from payload")
		return
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	msgId, err := h.manager.SendImage(r.Context(), name, req.Phone, data, req.Caption, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) SendVideo(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	data, err := decodeBase64Data(req.Video)
	if err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode base64 encoded data from payload")
		return
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "video/mp4"
	}

	msgId, err := h.manager.SendVideo(r.Context(), name, req.Phone, data, req.Caption, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) SendDocument(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	if req.FileName == "" {
		dto.Error(w, http.StatusBadRequest, "missing FileName in Payload")
		return
	}

	data, err := decodeBase64Data(req.Document)
	if err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode base64 encoded data from payload")
		return
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	msgId, err := h.manager.SendDocument(r.Context(), name, req.Phone, data, req.FileName, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) SendAudio(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendAudioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	data, err := decodeBase64Data(req.Audio)
	if err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode base64 encoded data from payload")
		return
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "audio/ogg; codecs=opus"
	}

	msgId, err := h.manager.SendAudio(r.Context(), name, req.Phone, data, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) SendSticker(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendStickerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	data, err := decodeBase64Data(req.Sticker)
	if err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode base64 encoded data from payload")
		return
	}

	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "image/webp"
	}

	msgId, err := h.manager.SendSticker(r.Context(), name, req.Phone, data, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) SendLocation(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	msgId, err := h.manager.SendLocation(r.Context(), name, req.Phone, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) SendContact(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	if req.Vcard == "" {
		dto.Error(w, http.StatusBadRequest, "missing Vcard in Payload")
		return
	}

	msgId, err := h.manager.SendContact(r.Context(), name, req.Phone, req.Name, req.Vcard)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) SendPoll(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendPollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	msgId, err := h.manager.SendPoll(r.Context(), name, req.Phone, req.Question, req.Options, req.MultiSelect)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) React(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.SendReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	emoji := req.Emoji
	if emoji == "remove" {
		emoji = ""
	}

	msgId, err := h.manager.SendReaction(r.Context(), name, req.Phone, req.MessageId, emoji)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	messageId := chi.URLParam(r, "messageId")

	var req dto.EditMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	msgId, err := h.manager.EditMessage(r.Context(), name, req.Phone, messageId, req.Body)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	messageId := chi.URLParam(r, "messageId")

	var req dto.RevokeMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	msgId, err := h.manager.RevokeMessage(r.Context(), name, req.Phone, messageId)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId})
}

func (h *MessageHandler) CheckPhone(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.CheckPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if len(req.Phone) == 0 {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	results, err := h.manager.CheckPhone(r.Context(), name, req.Phone)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	list := make([]dto.CheckPhoneResponse, 0, len(results))
	for _, res := range results {
		resp := dto.CheckPhoneResponse{Query: res.Query, IsInWhatsapp: res.IsIn}
		if res.IsIn {
			resp.JID = res.JID.String()
		}
		list = append(list, resp)
	}

	dto.Success(w, map[string]interface{}{"Users": list})
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

	devices := make([]string, len(userInfo.Devices))
	for i, d := range userInfo.Devices {
		devices[i] = d.String()
	}

	dto.Success(w, dto.UserInfoResponse{
		JID:       jid.String(),
		Status:    userInfo.Status,
		PictureID: userInfo.PictureID,
		Devices:   devices,
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

	dto.Success(w, map[string]interface{}{
		"URL":        pic.URL,
		"ID":         pic.ID,
		"Type":       pic.Type,
		"DirectPath": pic.DirectPath,
	})
}

func decodeBase64Data(data string) ([]byte, error) {
	if strings.Contains(data, ",") {
		parts := strings.SplitN(data, ",", 2)
		if len(parts) == 2 {
			data = parts[1]
		}
	}
	return base64.StdEncoding.DecodeString(data)
}
