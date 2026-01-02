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

// SignUp godoc
// @Summary      Sign up a new user
// @Description  Creates a new user with the provided details
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body dto.SignUpRequest true "Sign Up Request"
// @Success      201  {object}  dto.AuthResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /auth/signup [post]
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

// Login godoc
// @Summary      Logs in a user
// @Description  Logs in a user and returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body dto.LoginRequest true "Login Request"
// @Success      200  {object}  dto.AuthResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /auth/login [post]
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

// GetMyProfile godoc
// @Summary      Get user profile
// @Description  Get the profile of the currently logged-in user
// @Tags         Auth
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.AuthUser
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /me [get]
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
