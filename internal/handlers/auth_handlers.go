package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/falasefemi2/vendorhub/internal/dto"
	"github.com/falasefemi2/vendorhub/internal/service"
	"github.com/falasefemi2/vendorhub/internal/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: auth}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.StoreName == "" {
		utils.WriteError(w, http.StatusBadRequest, "store_name is required")
		return
	}
	user, err := h.authService.SignUp(req)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.Login(req)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, token)
}

func (h *AuthHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}
	user, err := h.authService.GetMyProfile(userID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}
