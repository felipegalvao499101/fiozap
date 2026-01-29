package handlers

import (
	"encoding/json"
	"net/http"

	"fiozap/internal/api/dto"
	"fiozap/internal/domain"

	"github.com/go-chi/chi/v5"
)

type ContactHandler struct {
	provider domain.Provider
}

func NewContactHandler(provider domain.Provider) *ContactHandler {
	return &ContactHandler{provider: provider}
}

// CheckPhone godoc
// @Summary      Verificar numeros
// @Description  Verifica se os numeros informados estao no WhatsApp
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.CheckPhoneRequest true "Lista de numeros"
// @Success      200 {object} dto.Response{data=dto.ContactsCheckResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/contacts/check [post]
func (h *ContactHandler) CheckPhone(w http.ResponseWriter, r *http.Request) {
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

	results, err := h.provider.CheckPhone(r.Context(), name, req.Phone)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	list := make([]dto.CheckPhoneResponse, 0, len(results))
	for _, res := range results {
		resp := dto.CheckPhoneResponse{Query: res.Phone, IsInWhatsapp: res.IsOnWhatsApp}
		if res.IsOnWhatsApp {
			resp.JID = res.JID
		}
		list = append(list, resp)
	}

	dto.Success(w, map[string]interface{}{"Contacts": list})
}

// GetInfo godoc
// @Summary      Obter info do contato
// @Description  Retorna informacoes de um contato (status, foto, devices)
// @Tags         contacts
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        phone path string true "Numero do telefone"
// @Success      200 {object} dto.Response{data=dto.UserInfoResponse}
// @Failure      404 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/contacts/{phone} [get]
func (h *ContactHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	phone := chi.URLParam(r, "phone")

	info, err := h.provider.GetUserInfo(r.Context(), name, []string{phone})
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	userInfo, ok := info[phone]
	if !ok {
		dto.Error(w, http.StatusNotFound, "contact not found")
		return
	}

	dto.Success(w, dto.UserInfoResponse{
		JID:       phone,
		Status:    userInfo.Status,
		PictureID: userInfo.PictureID,
		Devices:   userInfo.Devices,
	})
}

// GetAvatar godoc
// @Summary      Obter foto de perfil
// @Description  Retorna URL da foto de perfil do contato
// @Tags         contacts
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        phone path string true "Numero do telefone"
// @Success      200 {object} dto.Response{data=dto.AvatarResponse}
// @Failure      404 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/contacts/{phone}/avatar [get]
func (h *ContactHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	phone := chi.URLParam(r, "phone")

	pic, err := h.provider.GetProfilePicture(r.Context(), name, phone)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if pic == nil {
		dto.Error(w, http.StatusNotFound, "no profile picture")
		return
	}

	dto.Success(w, dto.AvatarResponse{
		URL: pic.URL,
		ID:  pic.ID,
	})
}

// GetBusinessProfile godoc
// @Summary      Obter perfil comercial
// @Description  Retorna perfil comercial do contato (nao implementado)
// @Tags         contacts
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        phone path string true "Numero do telefone"
// @Success      200 {object} dto.Response{data=dto.BusinessProfileResponse}
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/contacts/{phone}/business [get]
func (h *ContactHandler) GetBusinessProfile(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}
