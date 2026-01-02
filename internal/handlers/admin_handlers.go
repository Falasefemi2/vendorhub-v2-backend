package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/falasefemi2/vendorhub/internal/service"
	"github.com/falasefemi2/vendorhub/internal/utils"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(admin *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: admin}
}

// ApproveVendor godoc
// @Summary      Approve a vendor
// @Description  Approves a vendor with the given ID
// @Tags         Admin
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "Vendor ID"
// @Success      204
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /admin/vendors/{id}/approve [post]
func (h *AdminHandler) ApproveVendor(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())
	if rctx != nil {
		log.Println("Chi RoutePattern:", rctx.RoutePattern())
	}
	adminID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}
	vendorID := chi.URLParam(r, "id")
	err = h.adminService.ApproveVendor(adminID, vendorID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListPendingVendors godoc
// @Summary      List pending vendors
// @Description  Lists all vendors that are pending approval
// @Tags         Admin
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   models.User
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /admin/vendors/pending [get]
func (h *AdminHandler) ListPendingVendors(w http.ResponseWriter, r *http.Request) {
	adminID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}
	vendors, err := h.adminService.ListPendingVendors(adminID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, vendors)
}

// ListApprovedVendors godoc
// @Summary      List approved vendors
// @Description  Lists all vendors that have been approved
// @Tags         Admin
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   models.User
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /admin/vendors/approved [get]
func (h *AdminHandler) ListApprovedVendors(w http.ResponseWriter, r *http.Request) {
	adminID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}
	vendors, err := h.adminService.ListApprovedVendors(adminID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, vendors)
}
