package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"fiozap/internal/api/dto"
	"fiozap/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/skip2/go-qrcode"
)

type SessionHandler struct {
	provider domain.Provider
}

func NewSessionHandler(provider domain.Provider) *SessionHandler {
	return &SessionHandler{provider: provider}
}

// Create godoc
// @Summary      Criar sessao
// @Description  Cria uma nova sessao WhatsApp com o nome informado
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSessionRequest true "Nome da sessao"
// @Success      201 {object} dto.Response{data=dto.SessionResponse}
// @Failure      400 {object} dto.Response
// @Failure      409 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions [post]
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

	session, err := h.provider.CreateSession(r.Context(), req.Name)
	if err != nil {
		dto.Error(w, http.StatusConflict, err.Error())
		return
	}

	resp := sessionToDTO(session)
	resp.Token = session.GetToken()
	dto.Created(w, resp)
}

// List godoc
// @Summary      Listar sessoes
// @Description  Lista todas as sessoes WhatsApp criadas
// @Tags         sessions
// @Produce      json
// @Success      200 {object} dto.Response{data=[]dto.SessionResponse}
// @Security     ApiKeyAuth
// @Router       /sessions [get]
func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	sessions := h.provider.ListSessions()

	list := make([]dto.SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		list = append(list, sessionToDTO(s))
	}

	dto.Success(w, list)
}

// Get godoc
// @Summary      Obter sessao
// @Description  Retorna informacoes de uma sessao especifica
// @Tags         sessions
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=dto.SessionResponse}
// @Failure      404 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name} [get]
func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	session, err := h.provider.GetSession(name)
	if err != nil {
		dto.Error(w, http.StatusNotFound, err.Error())
		return
	}

	dto.Success(w, sessionToDTO(session))
}

// Connect godoc
// @Summary      Conectar sessao
// @Description  Inicia conexao da sessao com WhatsApp (gera QR code se necessario)
// @Tags         sessions
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=dto.SessionResponse}
// @Failure      404 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/connect [post]
func (h *SessionHandler) Connect(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	session, err := h.provider.Connect(r.Context(), name)
	if err != nil {
		// Se a sessão não existe, retorna 404
		if strings.Contains(err.Error(), "not found") {
			dto.Error(w, http.StatusNotFound, err.Error())
			return
		}
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, sessionToDTO(session))
}

// GetQR godoc
// @Summary      Obter QR code
// @Description  Retorna QR code para autenticacao. Use ?format=image para PNG
// @Tags         sessions
// @Produce      json
// @Produce      image/png
// @Param        name path string true "Nome da sessao"
// @Param        format query string false "Formato: image para PNG" Enums(image)
// @Success      200 {object} dto.Response{data=dto.QRResponse}
// @Failure      404 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/qr [get]
func (h *SessionHandler) GetQR(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	session, err := h.provider.GetSession(name)
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

	dto.Success(w, dto.QRResponse{QRCode: code})
}

// Disconnect godoc
// @Summary      Desconectar sessao
// @Description  Desconecta a sessao do WhatsApp (mantem dados)
// @Tags         sessions
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=dto.SessionResponse}
// @Failure      404 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/disconnect [post]
func (h *SessionHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.provider.Disconnect(name); err != nil {
		dto.Error(w, http.StatusNotFound, err.Error())
		return
	}

	session, _ := h.provider.GetSession(name)
	dto.Success(w, sessionToDTO(session))
}

// Delete godoc
// @Summary      Deletar sessao
// @Description  Remove sessao e todos os dados associados
// @Tags         sessions
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      404 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name} [delete]
func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.provider.DeleteSession(r.Context(), name); err != nil {
		dto.Error(w, http.StatusNotFound, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Session deleted"})
}

// Logout godoc
// @Summary      Logout da sessao
// @Description  Faz logout da sessao (requer novo QR code para reconectar)
// @Tags         sessions
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/logout [post]
func (h *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.provider.Logout(r.Context(), name); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Logged out"})
}

func sessionToDTO(s domain.Session) dto.SessionResponse {
	return dto.SessionResponse{
		Name:      s.GetName(),
		JID:       s.GetJID(),
		Phone:     s.GetPhone(),
		PushName:  s.GetPushName(),
		Connected: s.IsConnected(),
	}
}
