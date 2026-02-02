package handlers

import (
	"encoding/json"
	"net/http"

	"fiozap/internal/api/dto"
	"fiozap/internal/core"

	"github.com/go-chi/chi/v5"
)

type BlocklistHandler struct {
	provider core.Provider
}

func NewBlocklistHandler(provider core.Provider) *BlocklistHandler {
	return &BlocklistHandler{provider: provider}
}

// GetBlocklist godoc
// @Summary      Listar bloqueados
// @Description  Retorna lista de contatos bloqueados
// @Tags         blocklist
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=dto.BlocklistResponse}
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/blocklist [get]
func (h *BlocklistHandler) GetBlocklist(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	blocklist, err := h.provider.GetBlocklist(r.Context(), name)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.BlocklistResponse{JIDs: blocklist})
}

// Block godoc
// @Summary      Bloquear contato
// @Description  Adiciona contato a lista de bloqueados
// @Tags         blocklist
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.BlockRequest true "Numero a bloquear"
// @Success      200 {object} dto.Response{data=dto.BlockActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/blocklist/block [post]
func (h *BlocklistHandler) Block(w http.ResponseWriter, r *http.Request) {
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

	blocklist, err := h.provider.BlockContact(r.Context(), name, req.Phone)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.BlockActionResponse{Details: "Contact blocked", Blocklist: blocklist})
}

// Unblock godoc
// @Summary      Desbloquear contato
// @Description  Remove contato da lista de bloqueados
// @Tags         blocklist
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.BlockRequest true "Numero a desbloquear"
// @Success      200 {object} dto.Response{data=dto.BlockActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/blocklist/unblock [post]
func (h *BlocklistHandler) Unblock(w http.ResponseWriter, r *http.Request) {
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

	blocklist, err := h.provider.UnblockContact(r.Context(), name, req.Phone)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.BlockActionResponse{Details: "Contact unblocked", Blocklist: blocklist})
}
