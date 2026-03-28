package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User userResponse `json:"user"`
}

type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeMutationErrors(w, []FieldError{
			{Field: "email", Messages: []string{"Email is required"}},
			{Field: "password", Messages: []string{"Password is required"}},
		})
		return
	}

	user, token, err := h.authSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60, // 30 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	writeOK(w, authResponse{
		User: userResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeMutationErrors(w, []FieldError{
			{Field: "name", Messages: []string{"Name is required"}},
			{Field: "email", Messages: []string{"Email is required"}},
			{Field: "password", Messages: []string{"Password is required"}},
		})
		return
	}

	user, token, err := h.authSvc.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	writeCreated(w, authResponse{
		User: userResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		_ = h.authSvc.Logout(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	writeMutationOK(w)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	perms := middleware.GetPermissions(r.Context())

	writeOK(w, map[string]interface{}{
		"user": userResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
		"permissions": perms,
	})
}
