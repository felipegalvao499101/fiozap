package handlers

import (
	"encoding/json"
	"net/http"

	"fiozap/internal/api/dto"
	"fiozap/internal/core"

	"github.com/go-chi/chi/v5"
)

type GroupHandler struct {
	provider core.Provider
}

func NewGroupHandler(provider core.Provider) *GroupHandler {
	return &GroupHandler{provider: provider}
}

// Create godoc
// @Summary      Criar grupo
// @Description  Cria um novo grupo WhatsApp
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.CreateGroupRequest true "Dados do grupo"
// @Success      201 {object} dto.Response{data=dto.GroupResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups [post]
func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req dto.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Name == "" {
		dto.Error(w, http.StatusBadRequest, "missing Name in Payload")
		return
	}

	info, err := h.provider.CreateGroup(r.Context(), name, req.Name, req.Participants)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Created(w, groupToDTO(info))
}

// List godoc
// @Summary      Listar grupos
// @Description  Lista todos os grupos da sessao
// @Tags         groups
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Success      200 {object} dto.Response{data=[]dto.GroupResponse}
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups [get]
func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	groups, err := h.provider.GetGroups(r.Context(), name)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	list := make([]dto.GroupResponse, len(groups))
	for i, g := range groups {
		list[i] = groupToDTO(g)
	}

	dto.Success(w, map[string]interface{}{"Groups": list})
}

// Get godoc
// @Summary      Obter grupo
// @Description  Retorna informacoes de um grupo
// @Tags         groups
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Success      200 {object} dto.Response{data=dto.GroupResponse}
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid} [get]
func (h *GroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	info, err := h.provider.GetGroupInfo(r.Context(), name, groupJid)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, groupToDTO(info))
}

// SetName godoc
// @Summary      Alterar nome do grupo
// @Description  Altera o nome de um grupo
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.GroupNameRequest true "Novo nome"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/name [put]
func (h *GroupHandler) SetName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.provider.SetGroupName(r.Context(), name, groupJid, req.Name); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ActionResponse{Details: "Group Name set successfully"})
}

// SetTopic godoc
// @Summary      Alterar descricao do grupo
// @Description  Altera a descricao/topico de um grupo
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.GroupDescriptionRequest true "Nova descricao"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/topic [put]
func (h *GroupHandler) SetTopic(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.provider.SetGroupTopic(r.Context(), name, groupJid, req.Topic); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ActionResponse{Details: "Group Topic set successfully"})
}

// SetPhoto godoc
// @Summary      Alterar foto do grupo
// @Description  Altera a foto de um grupo
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.GroupPhotoRequest true "Imagem base64"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/photo [put]
func (h *GroupHandler) SetPhoto(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	pictureId, err := h.provider.SetGroupPhoto(r.Context(), name, groupJid, req.Image)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Group Photo set successfully", "PictureID": pictureId})
}

// Leave godoc
// @Summary      Sair do grupo
// @Description  Sai de um grupo
// @Tags         groups
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/leave [post]
func (h *GroupHandler) Leave(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	if err := h.provider.LeaveGroup(r.Context(), name, groupJid); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ActionResponse{Details: "Left group successfully"})
}

// GetInviteLink godoc
// @Summary      Obter link de convite
// @Description  Retorna link de convite do grupo
// @Tags         groups
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Success      200 {object} dto.Response{data=dto.InviteLinkResponse}
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/invite [get]
func (h *GroupHandler) GetInviteLink(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	link, err := h.provider.GetGroupInviteLink(r.Context(), name, groupJid, false)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.InviteLinkResponse{InviteLink: link})
}

// RevokeInviteLink godoc
// @Summary      Revogar link de convite
// @Description  Revoga e gera novo link de convite
// @Tags         groups
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Success      200 {object} dto.Response{data=dto.InviteLinkResponse}
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/invite/revoke [post]
func (h *GroupHandler) RevokeInviteLink(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	link, err := h.provider.GetGroupInviteLink(r.Context(), name, groupJid, true)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.InviteLinkResponse{InviteLink: link})
}

// Join godoc
// @Summary      Entrar no grupo
// @Description  Entra em um grupo usando codigo de convite
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.JoinGroupRequest true "Codigo do convite"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/join [post]
func (h *GroupHandler) Join(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req dto.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if req.Code == "" {
		dto.Error(w, http.StatusBadRequest, "missing Code in Payload")
		return
	}

	groupJid, err := h.provider.JoinGroupWithLink(r.Context(), name, req.Code)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"JID": groupJid})
}

// GetInviteInfo godoc
// @Summary      Info do convite
// @Description  Retorna informacoes do grupo pelo codigo de convite
// @Tags         groups
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        code path string true "Codigo do convite"
// @Success      200 {object} dto.Response{data=dto.GroupResponse}
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/invite/{code} [get]
func (h *GroupHandler) GetInviteInfo(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	code := chi.URLParam(r, "code")

	info, err := h.provider.GetGroupInfoFromLink(r.Context(), name, code)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, groupToDTO(info))
}

// AddParticipants godoc
// @Summary      Adicionar participantes
// @Description  Adiciona participantes ao grupo
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.ParticipantsRequest true "Lista de numeros"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/participants [post]
func (h *GroupHandler) AddParticipants(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.ParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.provider.AddParticipants(r.Context(), name, groupJid, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ActionResponse{Details: "Participants added successfully"})
}

// RemoveParticipants godoc
// @Summary      Remover participantes
// @Description  Remove participantes do grupo
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.ParticipantsRequest true "Lista de numeros"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/participants [delete]
func (h *GroupHandler) RemoveParticipants(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.ParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.provider.RemoveParticipants(r.Context(), name, groupJid, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ActionResponse{Details: "Participants removed successfully"})
}

// PromoteParticipants godoc
// @Summary      Promover participantes
// @Description  Promove participantes a administradores
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.ParticipantsRequest true "Lista de numeros"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/participants/promote [post]
func (h *GroupHandler) PromoteParticipants(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.ParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.provider.PromoteParticipants(r.Context(), name, groupJid, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ActionResponse{Details: "Participants promoted successfully"})
}

// DemoteParticipants godoc
// @Summary      Rebaixar participantes
// @Description  Rebaixa administradores a participantes comuns
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.ParticipantsRequest true "Lista de numeros"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/participants/demote [post]
func (h *GroupHandler) DemoteParticipants(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.ParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.provider.DemoteParticipants(r.Context(), name, groupJid, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ActionResponse{Details: "Participants demoted successfully"})
}

// SetAnnounce godoc
// @Summary      Modo somente admins
// @Description  Define se apenas admins podem enviar mensagens
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.GroupSettingRequest true "Valor"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/settings/announce [put]
func (h *GroupHandler) SetAnnounce(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.provider.SetGroupAnnounce(r.Context(), name, groupJid, req.Value); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ActionResponse{Details: "Group announce setting updated successfully"})
}

// SetLocked godoc
// @Summary      Bloquear edicao de info
// @Description  Define se apenas admins podem editar info do grupo
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.GroupSettingRequest true "Valor"
// @Success      200 {object} dto.Response{data=dto.ActionResponse}
// @Failure      400 {object} dto.Response
// @Failure      500 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/settings/locked [put]
func (h *GroupHandler) SetLocked(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.provider.SetGroupLocked(r.Context(), name, groupJid, req.Value); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.ActionResponse{Details: "Group locked setting updated successfully"})
}

// SetApproval godoc
// @Summary      Modo aprovacao
// @Description  Define se novos membros precisam de aprovacao (nao implementado)
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.GroupSettingRequest true "Valor"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/settings/approval [put]
func (h *GroupHandler) SetApproval(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// GetRequestParticipants godoc
// @Summary      Listar solicitacoes
// @Description  Lista solicitacoes de entrada no grupo (nao implementado)
// @Tags         groups
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/requests [get]
func (h *GroupHandler) GetRequestParticipants(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// ApproveRequestParticipants godoc
// @Summary      Aprovar solicitacoes
// @Description  Aprova solicitacoes de entrada no grupo (nao implementado)
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.ApproveParticipantsRequest true "Lista de numeros"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/requests/approve [post]
func (h *GroupHandler) ApproveRequestParticipants(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// RejectRequestParticipants godoc
// @Summary      Rejeitar solicitacoes
// @Description  Rejeita solicitacoes de entrada no grupo (nao implementado)
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.RejectParticipantsRequest true "Lista de numeros"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/requests/reject [post]
func (h *GroupHandler) RejectRequestParticipants(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// SetMemberAddMode godoc
// @Summary      Modo adicao de membros
// @Description  Define quem pode adicionar membros (nao implementado)
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        groupJid path string true "JID do grupo"
// @Param        request body dto.GroupMemberAddModeRequest true "Modo"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/groups/{groupJid}/settings/memberadd [put]
func (h *GroupHandler) SetMemberAddMode(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// LinkGroup godoc
// @Summary      Vincular grupo a comunidade
// @Description  Vincula um grupo a uma comunidade (nao implementado)
// @Tags         community
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.LinkGroupRequest true "Dados do vinculo"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/community/link [post]
func (h *GroupHandler) LinkGroup(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// UnlinkGroup godoc
// @Summary      Desvincular grupo de comunidade
// @Description  Desvincula um grupo de uma comunidade (nao implementado)
// @Tags         community
// @Accept       json
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        request body dto.LinkGroupRequest true "Dados do vinculo"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/community/unlink [post]
func (h *GroupHandler) UnlinkGroup(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// GetSubGroups godoc
// @Summary      Listar subgrupos
// @Description  Lista subgrupos de uma comunidade (nao implementado)
// @Tags         community
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        communityJid path string true "JID da comunidade"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/community/{communityJid}/subgroups [get]
func (h *GroupHandler) GetSubGroups(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

// GetLinkedParticipants godoc
// @Summary      Participantes da comunidade
// @Description  Lista participantes de uma comunidade (nao implementado)
// @Tags         community
// @Produce      json
// @Param        name path string true "Nome da sessao"
// @Param        communityJid path string true "JID da comunidade"
// @Failure      501 {object} dto.Response
// @Security     ApiKeyAuth
// @Router       /sessions/{name}/community/{communityJid}/participants [get]
func (h *GroupHandler) GetLinkedParticipants(w http.ResponseWriter, r *http.Request) {
	dto.Error(w, http.StatusNotImplemented, "not implemented")
}

func groupToDTO(g *core.GroupInfo) dto.GroupResponse {
	participants := make([]dto.ParticipantResponse, len(g.Participants))
	for i, p := range g.Participants {
		participants[i] = dto.ParticipantResponse{
			JID:          p.JID,
			IsAdmin:      p.IsAdmin,
			IsSuperAdmin: p.IsSuperAdmin,
		}
	}

	return dto.GroupResponse{
		JID:          g.JID,
		Name:         g.Name,
		Topic:        g.Topic,
		GroupCreated: g.GroupCreated.Format("2006-01-02T15:04:05Z"),
		OwnerJID:     g.OwnerJID,
		Participants: participants,
		IsAnnounce:   g.IsAnnounce,
		IsLocked:     g.IsLocked,
	}
}
