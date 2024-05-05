package handlers

import (
	"net/http"
	"time"

	"github.com/jmjp/go-rbac/internal/adapters/deliveries/tokens"
	"github.com/jmjp/go-rbac/internal/core/entities"
	"github.com/jmjp/go-rbac/internal/core/ports"
)

type AuthHandler struct {
	usecases ports.AuthUsecase
}

func NewAuthHandler(usecases ports.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecases: usecases,
	}
}

type LoginRequest struct {
	Email    string  `json:"email"`
	Avatar   *string `json:"avatar"`
	Username *string `json:"username"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body := new(LoginRequest)
	if err := getBody(r, body); err != nil {
		sendStatus(w, http.StatusBadRequest)
		return
	}
	msg, err := h.usecases.Login(body.Email, body.Avatar, body.Username)
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	sendString(w, http.StatusOK, *msg)
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")
	user, session, err := h.usecases.Verify(code, email, r.RemoteAddr, r.UserAgent())
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	accessToken, err := tokens.GeneratePasetoToken(user.ID.Hex(), user.Email, 5*time.Minute, user.Teams)
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	http.SetCookie(w, configureCookie("_refresh", session.Hash, session.ExpiresAt))
	sendJson(w, http.StatusOK, map[string]interface{}{"user": user, "access_token": accessToken})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session := r.URL.Query().Get("session")
	if session == "" {
		ses, err := r.Cookie("_refresh")
		if err != nil {
			sendString(w, http.StatusBadRequest, "refresh token missing")
			return
		}
		session = ses.Value
		http.SetCookie(w, configureCookie("_refresh", "", time.Now()))
	}
	userId := r.Context().Value("logged_user").(*tokens.Payload).UserId
	err := h.usecases.Logout(userId, session)
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	sendStatus(w, http.StatusOK)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_refresh")
	if err != nil {
		sendString(w, http.StatusBadRequest, "refresh token missing")
		return
	}
	user, session, err := h.usecases.Refresh(cookie.Value)
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	accessToken, err := tokens.GeneratePasetoToken(user.ID.Hex(), user.Email, 5*time.Minute, user.Teams)
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	http.SetCookie(w, configureCookie("_refresh", session.Hash, session.ExpiresAt))
	sendJson(w, http.StatusOK, map[string]interface{}{"user": user, "access_token": accessToken})
}

type SessionsResponse struct {
	*entities.Session
	Current bool `json:"current"`
}

func (h *AuthHandler) Sessions(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("logged_user").(*tokens.Payload).UserId
	cookie, err := r.Cookie("_refresh")
	if err != nil {
		sendString(w, http.StatusBadRequest, "refresh token missing")
		return
	}
	sessions, err := h.usecases.Sessions(userId)
	if err != nil {
		sendString(w, http.StatusBadRequest, err.Error())
		return
	}
	var output []*SessionsResponse
	for _, session := range sessions {
		output = append(output, &SessionsResponse{Session: session, Current: session.Hash == cookie.Value})
	}
	sendJson(w, http.StatusOK, output)
}
