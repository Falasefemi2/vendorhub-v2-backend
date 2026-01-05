package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/falasefemi2/vendorhub/internal/dto"
	"github.com/falasefemi2/vendorhub/internal/service"
	"github.com/falasefemi2/vendorhub/internal/utils"
)

type StoreHandler struct {
	userService    *service.AuthService
	productService *service.ProductService
}

func NewStoreHandler(userService *service.AuthService, productService *service.ProductService) *StoreHandler {
	return &StoreHandler{
		userService:    userService,
		productService: productService,
	}
}

// GetStoreBySlug godoc
// @Summary      Get store by slug
// @Description  Retrieves vendor's store and products by store slug (WhatsApp shareable link)
// @Tags         Stores
// @Accept       json
// @Produce      json
// @Param        slug path string true "Store slug (e.g., pizzahut-lagos)"
// @Success      200  {object}  dto.StoreDetailsResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /stores/{slug} [get]
func (sh *StoreHandler) GetStoreBySlug(w http.ResponseWriter, r *http.Request) {
	slugName := chi.URLParam(r, "slug")
	if slugName == "" {
		utils.WriteError(w, http.StatusBadRequest, "store slug is required")
		return
	}
	vendor, err := sh.userService.GetVendorBySlug(slugName)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}
	products, err := sh.productService.GetActiveProductsByUserID(r.Context(), vendor.ID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	response := &dto.StoreDetailsResponse{
		Store: &dto.StoreResponse{
			ID:             vendor.ID,
			Name:           vendor.StoreName,
			Slug:           vendor.StoreSlug,
			Username:       vendor.Username,
			Bio:            vendor.Bio,
			WhatsappNumber: vendor.WhatsappNumber,
			Email:          vendor.Email,
			CreatedAt:      vendor.CreatedAt.Format(time.RFC3339),
		},
		Products: products,
		StoreURL: "https://localhost:3000/stores/" + vendor.StoreSlug,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetStoreByVendorID godoc
// @Summary      Get store by vendor ID
// @Description  Retrieves vendor's store and products by vendor ID
// @Tags         Stores
// @Accept       json
// @Produce      json
// @Param        id query string true "Vendor ID"
// @Success      200  {object}  dto.StoreDetailsResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /stores/vendor [get]
func (sh *StoreHandler) GetStoreByVendorID(w http.ResponseWriter, r *http.Request) {
	vendorID := r.URL.Query().Get("id")
	if vendorID == "" {
		utils.WriteError(w, http.StatusBadRequest, "vendor id is required")
		return
	}

	// Get vendor
	vendor, err := sh.userService.GetUserByID(vendorID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	// Get vendor's active products
	products, err := sh.productService.GetActiveProductsByUserID(r.Context(), vendor.ID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	// Build store response
	response := &dto.StoreDetailsResponse{
		Store: &dto.StoreResponse{
			ID:             vendor.ID,
			Name:           vendor.StoreName,
			Slug:           vendor.StoreSlug,
			Username:       vendor.Username,
			Bio:            vendor.Bio,
			WhatsappNumber: vendor.WhatsappNumber,
			Email:          vendor.Email,
			CreatedAt:      vendor.CreatedAt.Format(time.RFC3339),
		},
		Products: products,
		StoreURL: "https://localhost:3000/stores/" + vendor.StoreSlug,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// UpdateMyStore godoc
// @Summary      Update authenticated vendor's store
// @Description  Updates the authenticated vendor's store information
// @Tags         Stores
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body body dto.UpdateStoreRequest true "Update Store Request"
// @Success      200  {object}  dto.StoreResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /stores/my [put]
func (sh *StoreHandler) UpdateMyStore(w http.ResponseWriter, r *http.Request) {
	vendorID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	role, err := utils.GetRoleFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	log.Printf("Role from context: %s", role)

	if role != "vendor" {
		utils.WriteError(w, http.StatusForbidden, "only vendors can update their store")
		return
	}

	var req dto.UpdateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	response, err := sh.userService.UpdateVendorStore(r.Context(), vendorID, req)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetMyStore godoc
// @Summary      Get authenticated vendor's store
// @Description  Retrieves the authenticated vendor's store and products
// @Tags         Stores
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.StoreDetailsResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      403  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /stores/my [get]
func (sh *StoreHandler) GetMyStore(w http.ResponseWriter, r *http.Request) {
	vendorID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	role, err := utils.GetRoleFromContext(r.Context())
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	if role != "vendor" {
		utils.WriteError(w, http.StatusForbidden, "only vendors can view their store")
		return
	}

	vendor, err := sh.userService.GetUserByID(vendorID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	products, err := sh.productService.GetProductsByUserID(r.Context(), vendorID)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	response := &dto.StoreDetailsResponse{
		Store: &dto.StoreResponse{
			ID:             vendor.ID,
			Name:           vendor.StoreName,
			Slug:           vendor.StoreSlug,
			Username:       vendor.Username,
			Bio:            vendor.Bio,
			WhatsappNumber: vendor.WhatsappNumber,
			Email:          vendor.Email,
			CreatedAt:      vendor.CreatedAt.Format(time.RFC3339),
		},
		Products: products,
		StoreURL: "https://localhost:3000/stores/" + vendor.StoreSlug,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetAllStores godoc
// @Summary      Get all vendor stores
// @Description  Retrieves all active vendor stores with pagination support
// @Tags         Stores
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number (default: 1)"
// @Param        page_size query int false "Page size (default: 20, max: 100)"
// @Success      200  {array}   dto.StoreResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /stores [get]
func (sh *StoreHandler) GetAllStores(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	pageSize := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	vendors, err := sh.userService.GetAllActiveVendors(page, pageSize)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	var storeResponses []*dto.StoreResponse
	for _, vendor := range vendors {
		storeResponses = append(storeResponses, &dto.StoreResponse{
			ID:             vendor.ID,
			Name:           vendor.StoreName,
			Slug:           vendor.StoreSlug,
			Username:       vendor.Username,
			Bio:            vendor.Bio,
			WhatsappNumber: vendor.WhatsappNumber,
			Email:          vendor.Email,
			CreatedAt:      vendor.CreatedAt.Format(time.RFC3339),
		})
	}

	utils.WriteJSON(w, http.StatusOK, storeResponses)
}

// SearchStores godoc
// @Summary      Search stores
// @Description  Searches for vendor stores by name or username
// @Tags         Stores
// @Accept       json
// @Produce      json
// @Param        q query string true "Search term"
// @Success      200  {array}   dto.StoreResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /stores/search [get]
func (sh *StoreHandler) SearchStores(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		utils.WriteError(w, http.StatusBadRequest, "search term is required")
		return
	}

	vendors, err := sh.userService.SearchVendors(searchTerm)
	if err != nil {
		utils.HandleServiceError(w, err)
		return
	}

	var storeResponses []*dto.StoreResponse
	for _, vendor := range vendors {
		storeResponses = append(storeResponses, &dto.StoreResponse{
			ID:             vendor.ID,
			Name:           vendor.StoreName,
			Slug:           vendor.StoreSlug,
			Username:       vendor.Username,
			Bio:            vendor.Bio,
			WhatsappNumber: vendor.WhatsappNumber,
			Email:          vendor.Email,
			CreatedAt:      vendor.CreatedAt.Format(time.RFC3339),
		})
	}

	utils.WriteJSON(w, http.StatusOK, storeResponses)
}
