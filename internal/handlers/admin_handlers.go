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
