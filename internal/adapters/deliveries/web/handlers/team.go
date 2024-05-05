package handlers

import (
	"net/http"

	"github.com/jmjp/go-rbac/internal/adapters/deliveries/tokens"
	"github.com/jmjp/go-rbac/internal/core/ports"
)

type TeamHandler struct {
	usecases ports.TeamUsecase
}

func NewTeamHandler(usecases ports.TeamUsecase) *TeamHandler {
	return &TeamHandler{
		usecases: usecases,
	}
}

type TeamRequest struct {
	Name string `json:"name"`
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	body := new(TeamRequest)
	if err := getBody(r, body); err != nil {
		sendStatus(w, http.StatusBadRequest)
		return
	}
	user := r.Context().Value("logged_user").(*tokens.Payload).UserId
	team, err := h.usecases.Create(user, body.Name)
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	sendJson(w, http.StatusOK, team)
}

func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	body := new(TeamRequest)
	if err := getBody(r, body); err != nil {
		sendStatus(w, http.StatusBadRequest)
		return
	}
	teamId := r.PathValue("teamId")
	user := r.Context().Value("logged_user").(*tokens.Payload).UserId
	team, err := h.usecases.Update(user, teamId, body.Name)
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	sendJson(w, http.StatusOK, team)
}

func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("teamId")
	user := r.Context().Value("logged_user").(*tokens.Payload).UserId
	err := h.usecases.Delete(user, teamId)
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	sendStatus(w, http.StatusOK)
}
