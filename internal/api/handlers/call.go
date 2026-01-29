package handlers

import (
	"encoding/json"
	"net/http"

	"fiozap/internal/api/dto"
	"fiozap/internal/domain"

	"github.com/go-chi/chi/v5"
)

type CallHandler struct {
	provider domain.Provider
}

func NewCallHandler(provider domain.Provider) *CallHandler {
	return &CallHandler{provider: provider}
}

// RejectCall godoc
// @Summary      Rejeitar chamada
// @Description  Rejeita uma chamada recebida
// @Tags         calls
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.RejectCallRequest true "Dados da chamada"
// @Success      200 {object} dto.Response{data=dto.CallActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/calls/reject [post]
func (h *CallHandler) RejectCall(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.RejectCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.CallFrom == "" || req.CallID == "" {
		dto.Error(w, http.StatusBadRequest, "missing CallFrom or CallID in Payload")
		return
	}

	if err := h.provider.RejectCall(r.Context(), name, req.CallFrom, req.CallID); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.CallActionResponse{Details: "Call rejected"})
}
