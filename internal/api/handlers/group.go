package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fiozap/fiozap/internal/api/dto"
	"github.com/fiozap/fiozap/internal/zap"
	"github.com/go-chi/chi/v5"
)

type GroupHandler struct {
	manager *zap.Manager
}

func NewGroupHandler(manager *zap.Manager) *GroupHandler {
	return &GroupHandler{manager: manager}
}

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

	info, err := h.manager.CreateGroup(r.Context(), req.Name, name, req.Participants)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Created(w, groupToDTO(info))
}

func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	groups, err := h.manager.GetGroups(r.Context(), name)
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

func (h *GroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	info, err := h.manager.GetGroupInfo(r.Context(), name, groupJid)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, groupToDTO(info))
}

func (h *GroupHandler) SetName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.manager.SetGroupName(r.Context(), name, groupJid, req.Name); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Group Name set successfully"})
}

func (h *GroupHandler) SetTopic(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.manager.SetGroupDescription(r.Context(), name, groupJid, req.Topic); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Group Topic set successfully"})
}

func (h *GroupHandler) SetPhoto(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	pictureId, err := h.manager.SetGroupPhoto(r.Context(), name, groupJid, req.Image)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Group Photo set successfully", "PictureID": pictureId})
}

func (h *GroupHandler) Leave(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	if err := h.manager.LeaveGroup(r.Context(), name, groupJid); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Left group successfully"})
}

func (h *GroupHandler) GetInviteLink(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	link, err := h.manager.GetGroupInviteLink(r.Context(), name, groupJid, false)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.InviteLinkResponse{InviteLink: link})
}

func (h *GroupHandler) RevokeInviteLink(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	link, err := h.manager.GetGroupInviteLink(r.Context(), name, groupJid, true)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, dto.InviteLinkResponse{InviteLink: link})
}

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

	groupJid, err := h.manager.JoinGroupWithLink(r.Context(), name, req.Code)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"JID": groupJid})
}

func (h *GroupHandler) GetInviteInfo(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	code := chi.URLParam(r, "code")

	info, err := h.manager.GetGroupInfoFromLink(r.Context(), name, code)
	if err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, groupToDTO(info))
}

func (h *GroupHandler) AddParticipants(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.ParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.manager.AddGroupParticipants(r.Context(), name, groupJid, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Participants added successfully"})
}

func (h *GroupHandler) RemoveParticipants(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.ParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.manager.RemoveGroupParticipants(r.Context(), name, groupJid, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Participants removed successfully"})
}

func (h *GroupHandler) PromoteParticipants(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.ParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.manager.PromoteGroupParticipants(r.Context(), name, groupJid, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Participants promoted successfully"})
}

func (h *GroupHandler) DemoteParticipants(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.ParticipantsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.manager.DemoteGroupParticipants(r.Context(), name, groupJid, req.Phone); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Participants demoted successfully"})
}

func (h *GroupHandler) SetAnnounce(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.manager.SetGroupAnnounce(r.Context(), name, groupJid, req.Value); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Group announce setting updated successfully"})
}

func (h *GroupHandler) SetLocked(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.manager.SetGroupLocked(r.Context(), name, groupJid, req.Value); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Group locked setting updated successfully"})
}

func (h *GroupHandler) SetApproval(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	groupJid := chi.URLParam(r, "groupJid")

	var req dto.GroupSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.Error(w, http.StatusBadRequest, "could not decode Payload")
		return
	}

	if err := h.manager.SetGroupJoinApproval(r.Context(), name, groupJid, req.Value); err != nil {
		dto.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(w, map[string]string{"Details": "Group approval setting updated successfully"})
}

func groupToDTO(g *zap.GroupInfo) dto.GroupResponse {
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
		GroupCreated: g.GroupCreated,
		OwnerJID:     g.OwnerJID,
		Participants: participants,
		IsAnnounce:   g.IsAnnounce,
		IsLocked:     g.IsLocked,
	}
}
