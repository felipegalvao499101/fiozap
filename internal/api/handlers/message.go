package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fiozap/internal/api/dto"
	"fiozap/internal/api/utils"
	"fiozap/internal/core"

	"github.com/go-chi/chi/v5"
)

type MessageHandler struct {
	provider core.Provider
}

func NewMessageHandler(provider core.Provider) *MessageHandler {
	return &MessageHandler{provider: provider}
}

// SendText godoc
// @Summary      Enviar texto
// @Description  Envia mensagem de texto para um contato ou grupo
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendTextRequest true "Dados da mensagem"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/text [post]
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

	msgId, err := h.provider.SendText(r.Context(), name, req.Phone, req.Body)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// SendImage godoc
// @Summary      Enviar imagem
// @Description  Envia imagem para um contato ou grupo. Aceita base64, data URL ou URL publica
// @Tags         messages
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendImageRequest true "Dados da imagem (JSON)"
// @Param        Phone formData string false "Numero do destinatario (form-data)"
// @Param        Caption formData string false "Legenda da imagem (form-data)"
// @Param        file formData file false "Arquivo de imagem (form-data)"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/image [post]
func (h *MessageHandler) SendImage(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var phone, caption, mimeType string
	var mediaData []byte

	contentType := r.Header.Get("Content-Type")

	// Verifica se é multipart/form-data
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(50 << 20); err != nil { // 50MB max
			dto.Error(w, http.StatusBadRequest, "failed to parse multipart form")
			return
		}

		phone = r.FormValue("Phone")
		caption = r.FormValue("Caption")

		media, err := utils.ProcessFormFile(r, "file")
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		mimeType = media.MimeType
	} else {
		// JSON request
		var req dto.SendImageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			dto.Error(w, http.StatusBadRequest, "could not decode Payload")
			return
		}

		phone = req.Phone
		caption = req.Caption
		mimeType = req.MimeType

		media, err := utils.ProcessMedia(req.Image, req.MimeType)
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		if mimeType == "" {
			mimeType = media.MimeType
		}
	}

	if phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	msgId, err := h.provider.SendImage(r.Context(), name, phone, mediaData, caption, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// SendVideo godoc
// @Summary      Enviar video
// @Description  Envia video para um contato ou grupo. Aceita base64, data URL ou URL publica
// @Tags         messages
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendVideoRequest true "Dados do video (JSON)"
// @Param        Phone formData string false "Numero do destinatario (form-data)"
// @Param        Caption formData string false "Legenda do video (form-data)"
// @Param        file formData file false "Arquivo de video (form-data)"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/video [post]
func (h *MessageHandler) SendVideo(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var phone, caption, mimeType string
	var mediaData []byte

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(100 << 20); err != nil { // 100MB max
			dto.Error(w, http.StatusBadRequest, "failed to parse multipart form")
			return
		}

		phone = r.FormValue("Phone")
		caption = r.FormValue("Caption")

		media, err := utils.ProcessFormFile(r, "file")
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		mimeType = media.MimeType
	} else {
		var req dto.SendVideoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			dto.Error(w, http.StatusBadRequest, "could not decode Payload")
			return
		}

		phone = req.Phone
		caption = req.Caption
		mimeType = req.MimeType

		media, err := utils.ProcessMedia(req.Video, req.MimeType)
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		if mimeType == "" {
			mimeType = media.MimeType
		}
	}

	if phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	if mimeType == "" {
		mimeType = "video/mp4"
	}

	msgId, err := h.provider.SendVideo(r.Context(), name, phone, mediaData, caption, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// SendDocument godoc
// @Summary      Enviar documento
// @Description  Envia documento para um contato ou grupo. Aceita base64, data URL ou URL publica
// @Tags         messages
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendDocumentRequest true "Dados do documento (JSON)"
// @Param        Phone formData string false "Numero do destinatario (form-data)"
// @Param        FileName formData string false "Nome do arquivo (form-data)"
// @Param        file formData file false "Arquivo (form-data)"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/document [post]
func (h *MessageHandler) SendDocument(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var phone, fileName, mimeType string
	var mediaData []byte

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(100 << 20); err != nil { // 100MB max
			dto.Error(w, http.StatusBadRequest, "failed to parse multipart form")
			return
		}

		phone = r.FormValue("Phone")
		fileName = r.FormValue("FileName")

		media, err := utils.ProcessFormFile(r, "file")
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		mimeType = media.MimeType
		if fileName == "" {
			fileName = media.FileName
		}
	} else {
		var req dto.SendDocumentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			dto.Error(w, http.StatusBadRequest, "could not decode Payload")
			return
		}

		phone = req.Phone
		fileName = req.FileName
		mimeType = req.MimeType

		media, err := utils.ProcessMedia(req.Document, req.MimeType)
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		if mimeType == "" {
			mimeType = media.MimeType
		}
		if fileName == "" && media.FileName != "" {
			fileName = media.FileName
		}
	}

	if phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	if fileName == "" {
		dto.Error(w, http.StatusBadRequest, "missing FileName in Payload")
		return
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	msgId, err := h.provider.SendDocument(r.Context(), name, phone, mediaData, fileName, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// SendAudio godoc
// @Summary      Enviar audio
// @Description  Envia audio para um contato ou grupo. Aceita base64, data URL ou URL publica
// @Tags         messages
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendAudioRequest true "Dados do audio (JSON)"
// @Param        Phone formData string false "Numero do destinatario (form-data)"
// @Param        file formData file false "Arquivo de audio (form-data)"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/audio [post]
func (h *MessageHandler) SendAudio(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var phone, mimeType string
	var mediaData []byte

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(50 << 20); err != nil { // 50MB max
			dto.Error(w, http.StatusBadRequest, "failed to parse multipart form")
			return
		}

		phone = r.FormValue("Phone")

		media, err := utils.ProcessFormFile(r, "file")
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		mimeType = media.MimeType
	} else {
		var req dto.SendAudioRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			dto.Error(w, http.StatusBadRequest, "could not decode Payload")
			return
		}

		phone = req.Phone
		mimeType = req.MimeType

		media, err := utils.ProcessMedia(req.Audio, req.MimeType)
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		if mimeType == "" {
			mimeType = media.MimeType
		}
	}

	if phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	if mimeType == "" {
		mimeType = "audio/ogg; codecs=opus"
	}

	msgId, err := h.provider.SendAudio(r.Context(), name, phone, mediaData, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// SendSticker godoc
// @Summary      Enviar sticker
// @Description  Envia sticker para um contato ou grupo. Aceita base64, data URL ou URL publica
// @Tags         messages
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendStickerRequest true "Dados do sticker (JSON)"
// @Param        Phone formData string false "Numero do destinatario (form-data)"
// @Param        file formData file false "Arquivo de sticker (form-data)"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/sticker [post]
func (h *MessageHandler) SendSticker(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var phone, mimeType string
	var mediaData []byte

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
			dto.Error(w, http.StatusBadRequest, "failed to parse multipart form")
			return
		}

		phone = r.FormValue("Phone")

		media, err := utils.ProcessFormFile(r, "file")
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		mimeType = media.MimeType
	} else {
		var req dto.SendStickerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			dto.Error(w, http.StatusBadRequest, "could not decode Payload")
			return
		}

		phone = req.Phone
		mimeType = req.MimeType

		media, err := utils.ProcessMedia(req.Sticker, req.MimeType)
		if err != nil {
			dto.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		mediaData = media.Data
		if mimeType == "" {
			mimeType = media.MimeType
		}
	}

	if phone == "" {
		dto.Error(w, http.StatusBadRequest, "missing Phone in Payload")
		return
	}

	if mimeType == "" {
		mimeType = "image/webp"
	}

	msgId, err := h.provider.SendSticker(r.Context(), name, phone, mediaData, mimeType)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// SendLocation godoc
// @Summary      Enviar localizacao
// @Description  Envia localizacao para um contato ou grupo
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendLocationRequest true "Dados da localizacao"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/location [post]
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

	msgId, err := h.provider.SendLocation(r.Context(), name, req.Phone, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// SendContact godoc
// @Summary      Enviar contato
// @Description  Envia contato (vCard) para um contato ou grupo
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendContactRequest true "Dados do contato"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/contact [post]
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

	msgId, err := h.provider.SendContact(r.Context(), name, req.Phone, req.Name, req.Vcard)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// SendPoll godoc
// @Summary      Enviar enquete
// @Description  Envia enquete para um contato ou grupo
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendPollRequest true "Dados da enquete"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/poll [post]
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

	msgId, err := h.provider.SendPoll(r.Context(), name, req.Phone, req.Question, req.Options, req.MultiSelect)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// React godoc
// @Summary      Enviar reacao
// @Description  Envia reacao (emoji) para uma mensagem
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SendReactionRequest true "Dados da reacao"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/reaction [post]
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
	if emoji == "" || emoji == "remove" {
		emoji = ""
	}

	msgId, err := h.provider.SendReaction(r.Context(), name, req.Phone, req.MessageId, emoji)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// Edit godoc
// @Summary      Editar mensagem
// @Description  Edita uma mensagem enviada
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        messageId path string true "ID da mensagem"
// @Param        request body dto.EditMessageRequest true "Novo conteudo"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/{messageId} [put]
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

	msgId, err := h.provider.EditMessage(r.Context(), name, req.Phone, messageId, req.Body)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}

// Revoke godoc
// @Summary      Revogar mensagem
// @Description  Revoga/deleta uma mensagem enviada
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        messageId path string true "ID da mensagem"
// @Param        request body dto.RevokeMessageRequest true "Dados do contato"
// @Success      200 {object} dto.Response{data=dto.MessageResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/messages/{messageId} [delete]
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

	msgId, err := h.provider.RevokeMessage(r.Context(), name, req.Phone, messageId)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.MessageResponse{MessageId: msgId.ID})
}
