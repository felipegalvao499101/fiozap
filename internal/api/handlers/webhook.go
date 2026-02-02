package handlers

import (
	"encoding/json"
	"net/http"

	"fiozap/internal/api/dto"
	"fiozap/internal/integrations/webhook"

	"github.com/go-chi/chi/v5"
)

type WebhookHandler struct {
	dispatcher *webhook.Dispatcher
}

func NewWebhookHandler(dispatcher *webhook.Dispatcher) *WebhookHandler {
	return &WebhookHandler{dispatcher: dispatcher}
}

// SetWebhook godoc
// @Summary      Configurar webhook
// @Description  Configura a URL e eventos do webhook para a sessao
// @Tags         webhook
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SetWebhookRequest true "Configuracao do webhook"
// @Success      200 {object} dto.Response{data=dto.WebhookConfigResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/webhook [post]
func (h *WebhookHandler) SetWebhook(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.SetWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.WebhookURL == "" {
		dto.Error(w, http.StatusBadRequest, "missing WebhookURL in Payload")
		return
	}

	events := webhook.ParseEventTypes(req.Events)
	if len(events) == 0 {
		events = []webhook.EventType{webhook.EventAll}
	}

	h.dispatcher.SetConfig(name, req.WebhookURL, events)

	config := h.dispatcher.GetConfig(name)
	dto.Success(w, dto.WebhookConfigResponse{
		Webhook:    config.URL,
		Events:     webhook.EventTypesToStrings(config.Events),
		HMACKeySet: config.HMACKeySet,
	})
}

// GetWebhook godoc
// @Summary      Obter configuracao do webhook
// @Description  Retorna a configuracao atual do webhook para a sessao
// @Tags         webhook
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=dto.WebhookConfigResponse}
// @Failure      404 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/webhook [get]
func (h *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	config := h.dispatcher.GetConfig(name)
	if config == nil {
		dto.Success(w, dto.WebhookConfigResponse{
			Webhook: "",
			Events:  []string{},
		})
		return
	}

	dto.Success(w, dto.WebhookConfigResponse{
		Webhook:    config.URL,
		Events:     webhook.EventTypesToStrings(config.Events),
		HMACKeySet: config.HMACKeySet,
	})
}

// DeleteWebhook godoc
// @Summary      Remover webhook
// @Description  Remove a configuracao do webhook para a sessao
// @Tags         webhook
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=dto.WebhookActionResponse}
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/webhook [delete]
func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	h.dispatcher.RemoveConfig(name)
	dto.Success(w, dto.WebhookActionResponse{Details: "Webhook removed"})
}

// SetHMAC godoc
// @Summary      Configurar HMAC
// @Description  Configura a chave HMAC para assinatura dos webhooks
// @Tags         webhook
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.SetHMACRequest true "Chave HMAC"
// @Success      200 {object} dto.Response{data=dto.WebhookActionResponse}
// @Failure      400 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/webhook/hmac [post]
func (h *WebhookHandler) SetHMAC(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.SetHMACRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.dispatcher.SetHMACKey(name, req.HMACKey); err != nil {
		dto.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	dto.Success(w, dto.WebhookActionResponse{Details: "HMAC key configured"})
}

// DeleteHMAC godoc
// @Summary      Remover HMAC
// @Description  Remove a chave HMAC da sessao
// @Tags         webhook
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=dto.WebhookActionResponse}
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/webhook/hmac [delete]
func (h *WebhookHandler) DeleteHMAC(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	h.dispatcher.RemoveHMACKey(name)
	dto.Success(w, dto.WebhookActionResponse{Details: "HMAC key removed"})
}

// GetSupportedEvents godoc
// @Summary      Listar eventos suportados
// @Description  Retorna a lista de tipos de eventos suportados para webhook
// @Tags         webhook
// @Produce      json
// @Success      200 {object} dto.Response{data=dto.SupportedEventsResponse}
// @Security     ApiKeyAuth
// @Router       /webhook/events [get]
func (h *WebhookHandler) GetSupportedEvents(w http.ResponseWriter, r *http.Request) {
	dto.Success(w, dto.SupportedEventsResponse{
		Events: webhook.SupportedEventStrings(),
	})
}
